/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package startcmd

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/tink/go/subtle/random"
	ariescouchdbstorage "github.com/hyperledger/aries-framework-go-ext/component/storage/couchdb"
	ariesmongodbstorage "github.com/hyperledger/aries-framework-go-ext/component/storage/mongodb"
	ariesmysqlstorage "github.com/hyperledger/aries-framework-go-ext/component/storage/mysql"
	"github.com/hyperledger/aries-framework-go-ext/component/vdr/orb"
	ariesmemstorage "github.com/hyperledger/aries-framework-go/component/storageutil/mem"
	"github.com/hyperledger/aries-framework-go/pkg/crypto/tinkcrypto"
	ariesld "github.com/hyperledger/aries-framework-go/pkg/doc/ld"
	"github.com/hyperledger/aries-framework-go/pkg/doc/ldcontext/remote"
	vdrapi "github.com/hyperledger/aries-framework-go/pkg/framework/aries/api/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/hyperledger/aries-framework-go/pkg/kms/localkms"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock"
	"github.com/hyperledger/aries-framework-go/pkg/secretlock/local"
	vdrpkg "github.com/hyperledger/aries-framework-go/pkg/vdr"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/httpbinding"
	"github.com/hyperledger/aries-framework-go/pkg/vdr/key"
	ariesstorage "github.com/hyperledger/aries-framework-go/spi/storage"
	jsonld "github.com/piprate/json-gold/ld"
	tlsutils "github.com/trustbloc/edge-core/pkg/utils/tls"

	"github.com/trustbloc/vcs/pkg/ld"
	vcsstorage "github.com/trustbloc/vcs/pkg/storage"
	ariesvcsprovider "github.com/trustbloc/vcs/pkg/storage/ariesprovider"
	mongodbvcsprovider "github.com/trustbloc/vcs/pkg/storage/mongodbprovider"
)

// mode in which to run the vc-rest service
type mode string

const (
	verifier mode = "verifier"
	issuer   mode = "issuer"
	holder   mode = "holder"
	combined mode = "combined"
)

// Configuration for the vc-rest API server.
type Configuration struct {
	RootCAs           *x509.CertPool
	Storage           *vcStorageProviders
	LocalKMS          *localkms.LocalKMS
	Crypto            *tinkcrypto.Crypto
	VDR               vdrapi.Registry
	DocumentLoader    jsonld.DocumentLoader
	LDContextStore    *ld.StoreProvider
	StartupParameters *startupParameters
}

func prepareConfiguration(parameters *startupParameters) (*Configuration, error) {
	rootCAs, err := tlsutils.GetCertPool(parameters.tlsSystemCertPool, parameters.tlsCACerts)
	if err != nil {
		return nil, err
	}

	storeProviders, err := createStoreProviders(parameters)
	if err != nil {
		return nil, err
	}

	localKMS, err := createKMS(storeProviders.kmsSecretsProvider)
	if err != nil {
		return nil, err
	}

	crypto, err := tinkcrypto.New()
	if err != nil {
		return nil, err
	}

	vdr, err := createVDRI(parameters.universalResolverURL,
		&tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12}, parameters.blocDomain,
		parameters.requestTokens["sidetreeToken"])
	if err != nil {
		return nil, err
	}

	ldStore, err := ld.NewStoreProvider(storeProviders.provider)
	if err != nil {
		return nil, err
	}

	loader, err := createJSONLDDocumentLoader(ldStore, rootCAs, parameters.contextProviderURLs,
		parameters.contextEnableRemote)
	if err != nil {
		return nil, err
	}

	return &Configuration{
		RootCAs:           rootCAs,
		Storage:           storeProviders,
		LocalKMS:          localKMS,
		Crypto:            crypto,
		VDR:               vdr,
		DocumentLoader:    loader,
		LDContextStore:    ldStore,
		StartupParameters: parameters,
	}, nil
}

type vcStorageProviders struct {
	provider           vcsstorage.Provider
	kmsSecretsProvider vcsstorage.Provider
}

func createStoreProviders(parameters *startupParameters) (*vcStorageProviders, error) {
	var edgeServiceProvs vcStorageProviders

	var err error

	edgeServiceProvs.provider, err = createMainStoreProvider(parameters)
	if err != nil {
		return nil, err
	}

	edgeServiceProvs.kmsSecretsProvider, err = createKMSSecretsProvider(parameters)
	if err != nil {
		return nil, err
	}

	return &edgeServiceProvs, nil
}

func createMainStoreProvider(parameters *startupParameters) (vcsstorage.Provider, error) { //nolint: dupl
	switch {
	case strings.EqualFold(parameters.dbParameters.databaseType, databaseTypeMemOption):
		return ariesvcsprovider.New(ariesmemstorage.NewProvider()), nil
	case strings.EqualFold(parameters.dbParameters.databaseType, databaseTypeCouchDBOption):
		couchDBProvider, err := ariescouchdbstorage.NewProvider(parameters.dbParameters.databaseURL,
			ariescouchdbstorage.WithDBPrefix(parameters.dbParameters.databasePrefix))
		if err != nil {
			return nil, err
		}

		return ariesvcsprovider.New(couchDBProvider), nil
	case strings.EqualFold(parameters.dbParameters.databaseType, databaseTypeMYSQLDBOption):
		mySQLProvider, err := ariesmysqlstorage.NewProvider(parameters.dbParameters.databaseURL,
			ariesmysqlstorage.WithDBPrefix(parameters.dbParameters.databasePrefix))
		if err != nil {
			return nil, err
		}

		return ariesvcsprovider.New(mySQLProvider), nil
	case strings.EqualFold(parameters.dbParameters.databaseType, databaseTypeMongoDBOption):
		mongoDBProvider, err := ariesmongodbstorage.NewProvider(parameters.dbParameters.databaseURL,
			ariesmongodbstorage.WithDBPrefix(parameters.dbParameters.databasePrefix))
		if err != nil {
			return nil, err
		}

		return mongodbvcsprovider.New(mongoDBProvider), nil
	default:
		return nil, fmt.Errorf("%s is not a valid database type."+
			" run start --help to see the available options", parameters.dbParameters.databaseType)
	}
}

func createKMSSecretsProvider(parameters *startupParameters) (vcsstorage.Provider, error) { //nolint: dupl
	switch {
	case strings.EqualFold(parameters.dbParameters.kmsSecretsDatabaseType, databaseTypeMemOption):
		return ariesvcsprovider.New(ariesmemstorage.NewProvider()), nil
	case strings.EqualFold(parameters.dbParameters.kmsSecretsDatabaseType, databaseTypeCouchDBOption):
		couchDBProvider, err := ariescouchdbstorage.NewProvider(parameters.dbParameters.kmsSecretsDatabaseURL,
			ariescouchdbstorage.WithDBPrefix(parameters.dbParameters.kmsSecretsDatabasePrefix))
		if err != nil {
			return nil, err
		}

		return ariesvcsprovider.New(couchDBProvider), nil
	case strings.EqualFold(parameters.dbParameters.kmsSecretsDatabaseType, databaseTypeMYSQLDBOption):
		mySQLProvider, err := ariesmysqlstorage.NewProvider(parameters.dbParameters.kmsSecretsDatabaseURL,
			ariesmysqlstorage.WithDBPrefix(parameters.dbParameters.kmsSecretsDatabasePrefix))
		if err != nil {
			return nil, err
		}

		return ariesvcsprovider.New(mySQLProvider), nil
	case strings.EqualFold(parameters.dbParameters.kmsSecretsDatabaseType, databaseTypeMongoDBOption):
		mongoDBProvider, err := ariesmongodbstorage.NewProvider(parameters.dbParameters.kmsSecretsDatabaseURL,
			ariesmongodbstorage.WithDBPrefix(parameters.dbParameters.kmsSecretsDatabasePrefix))
		if err != nil {
			return nil, err
		}

		return mongodbvcsprovider.New(mongoDBProvider), nil
	default:
		return nil, fmt.Errorf("%s is not a valid KMS secrets database type."+
			" run start --help to see the available options", parameters.dbParameters.kmsSecretsDatabaseType)
	}
}

type kmsProvider struct {
	storageProvider   kms.Store
	secretLockService secretlock.Service
}

func (k kmsProvider) StorageProvider() kms.Store {
	return k.storageProvider
}

func (k kmsProvider) SecretLock() secretlock.Service {
	return k.secretLockService
}

func createKMS(kmsSecretsProvider vcsstorage.Provider) (*localkms.LocalKMS, error) {
	localKMS, err := createLocalKMS(kmsSecretsProvider)
	if err != nil {
		return nil, err
	}

	return localKMS, nil
}

func createLocalKMS(kmsSecretsStoreProvider vcsstorage.Provider) (*localkms.LocalKMS, error) {
	masterKeyReader, err := prepareMasterKeyReader(kmsSecretsStoreProvider)
	if err != nil {
		return nil, err
	}

	secretLockService, err := local.NewService(masterKeyReader, nil)
	if err != nil {
		return nil, err
	}

	// TODO (#769): Create our own implementation of the KMS storage interface and pass it in here instead of wrapping
	//  the Aries storage provider.
	kmsStore, err := kms.NewAriesProviderWrapper(kmsSecretsStoreProvider.GetAriesProvider())
	if err != nil {
		return nil, err
	}

	kmsProv := kmsProvider{
		storageProvider:   kmsStore,
		secretLockService: secretLockService,
	}

	return localkms.New(masterKeyURI, kmsProv)
}

// prepareMasterKeyReader prepares a master key reader for secret lock usage
func prepareMasterKeyReader(kmsSecretsStoreProvider vcsstorage.Provider) (*bytes.Reader, error) {
	masterKeyStore, err := kmsSecretsStoreProvider.OpenMasterKeyStore()
	if err != nil {
		return nil, err
	}

	masterKey, err := masterKeyStore.Get()
	if err != nil {
		if errors.Is(err, ariesstorage.ErrDataNotFound) {
			masterKey = random.GetRandomBytes(uint32(masterKeyNumBytes))

			putErr := masterKeyStore.Put(masterKey)
			if putErr != nil {
				return nil, putErr
			}
		} else {
			return nil, err
		}
	}

	masterKeyReader := bytes.NewReader(masterKey)

	return masterKeyReader, nil
}

func createVDRI(universalResolver string, tlsConfig *tls.Config, blocDomain,
	sidetreeAuthToken string) (vdrapi.Registry, error) {
	var opts []vdrpkg.Option

	if universalResolver != "" {
		universalResolverVDRI, err := httpbinding.New(universalResolver,
			httpbinding.WithAccept(acceptsDID), httpbinding.WithHTTPClient(&http.Client{
				Transport: &http.Transport{
					TLSClientConfig: tlsConfig,
				},
			}))
		if err != nil {
			return nil, fmt.Errorf("failed to create new universal resolver vdr: %w", err)
		}

		// add universal resolver vdr
		opts = append(opts, vdrpkg.WithVDR(universalResolverVDRI))
	}

	vdr, err := orb.New(nil, orb.WithDomain(blocDomain), orb.WithTLSConfig(tlsConfig),
		orb.WithAuthToken(sidetreeAuthToken))
	if err != nil {
		return nil, err
	}

	// add bloc vdr
	opts = append(opts, vdrpkg.WithVDR(vdr), vdrpkg.WithVDR(key.New()))

	return vdrpkg.New(opts...), nil
}

// acceptsDID returns if given did method is accepted by VC REST api
func acceptsDID(method string) bool {
	return method == didMethodVeres || method == didMethodElement || method == didMethodSov ||
		method == didMethodWeb || method == didMethodFactom
}

func createJSONLDDocumentLoader(ldStore *ld.StoreProvider, rootCAs *x509.CertPool,
	providerURLs []string, contextEnableRemote bool) (jsonld.DocumentLoader, error) {
	var loaderOpts []ariesld.DocumentLoaderOpts

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12},
		},
	}

	for _, url := range providerURLs {
		loaderOpts = append(loaderOpts,
			ariesld.WithRemoteProvider(
				remote.NewProvider(url, remote.WithHTTPClient(httpClient)),
			),
		)
	}

	if contextEnableRemote {
		loaderOpts = append(loaderOpts,
			ariesld.WithRemoteDocumentLoader(jsonld.NewDefaultDocumentLoader(http.DefaultClient)))
	}

	loader, err := ld.NewDocumentLoader(ldStore, loaderOpts...)
	if err != nil {
		return nil, err
	}

	return loader, nil
}
