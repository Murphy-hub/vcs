/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

//go:generate oapi-codegen --config=openapi.cfg.yaml ../../../../docs/v1/openapi.yaml
//go:generate mockgen -destination controller_mocks_test.go -self_package mocks -package oidc4ci_test . StateStore,OAuth2Provider,IssuerInteractionClient,HTTPClient,OAuth2Client

package oidc4ci

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hyperledger/aries-framework-go/pkg/doc/jose"
	"github.com/hyperledger/aries-framework-go/pkg/doc/jwt"
	"github.com/labstack/echo/v4"
	"github.com/ory/fosite"
	"github.com/samber/lo"
	"golang.org/x/oauth2"

	"github.com/trustbloc/vcs/pkg/oauth2client"
	"github.com/trustbloc/vcs/pkg/restapi/resterr"
	"github.com/trustbloc/vcs/pkg/restapi/v1/common"
	"github.com/trustbloc/vcs/pkg/restapi/v1/issuer"
	apiUtil "github.com/trustbloc/vcs/pkg/restapi/v1/util"
	"github.com/trustbloc/vcs/pkg/service/oidc4ci"
	"github.com/trustbloc/vcs/pkg/storage/mongodb/oidc4cistatestore"
)

const (
	sessionOpStateKey          = "opState"
	authorizationDetailsKey    = "authDetails"
	txIDKey                    = "txID"
	preAuthorizedCodeGrantType = "urn:ietf:params:oauth:grant-type:pre-authorized_code"
	cNonceKey                  = "cNonce"
	cNonceExpiresAtKey         = "cNonceExpiresAt"
	cNonceTTL                  = 5 * time.Minute
	cNonceSize                 = 15
)

// StateStore stores authorization request/response state.
type StateStore interface {
	SaveAuthorizeState(
		ctx context.Context,
		opState string,
		state *oidc4cistatestore.AuthorizeState,
		params ...func(insertOptions *oidc4ci.InsertOptions),
	) error

	GetAuthorizeState(
		ctx context.Context,
		opState string,
	) (*oidc4cistatestore.AuthorizeState, error)
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type OAuth2Client interface {
	GeneratePKCE() (verifier string, challenge string, method string, err error)
	AuthCodeURL(_ context.Context, cfg oauth2.Config, state string, opts ...oauth2client.AuthCodeOption) string
	Exchange(
		ctx context.Context,
		cfg oauth2.Config,
		code string,
		client *http.Client,
		opts ...oauth2client.AuthCodeOption,
	) (*oauth2.Token, error)
	AuthCodeURLWithPAR(
		ctx context.Context,
		cfg oauth2.Config,
		parEndpoint string,
		state string,
		client *http.Client,
		opts ...oauth2client.AuthCodeOption,
	) (string, error)
}

// OAuth2Provider provides functionality for OAuth2 handlers.
type OAuth2Provider fosite.OAuth2Provider

// IssuerInteractionClient defines API client for interaction with issuer private API.
type IssuerInteractionClient issuer.ClientInterface

// Config holds configuration options for Controller.
type Config struct {
	OAuth2Provider          OAuth2Provider
	StateStore              StateStore
	IssuerInteractionClient IssuerInteractionClient
	IssuerVCSPublicHost     string
	OAuth2Client            OAuth2Client
	PreAuthorizeClient      HTTPClient
	DefaultHTTPClient       *http.Client
	JWTVerifier             jose.SignatureVerifier
	ExternalHostURL         string
}

// Controller for OIDC credential issuance API.
type Controller struct {
	oauth2Provider          OAuth2Provider
	stateStore              StateStore
	issuerInteractionClient IssuerInteractionClient
	issuerVCSPublicHost     string
	oAuth2Client            OAuth2Client
	preAuthorizeClient      HTTPClient
	defaultHTTPClient       *http.Client
	jwtVerifier             jose.SignatureVerifier
	internalHostURL         string
}

// NewController creates a new Controller instance.
func NewController(config *Config) *Controller {
	return &Controller{
		oauth2Provider:          config.OAuth2Provider,
		stateStore:              config.StateStore,
		issuerInteractionClient: config.IssuerInteractionClient,
		issuerVCSPublicHost:     config.IssuerVCSPublicHost,
		oAuth2Client:            config.OAuth2Client,
		preAuthorizeClient:      config.PreAuthorizeClient,
		defaultHTTPClient:       config.DefaultHTTPClient,
		jwtVerifier:             config.JWTVerifier,
		internalHostURL:         config.ExternalHostURL,
	}
}

// OidcPushedAuthorizationRequest handles OIDC pushed authorization request (POST /oidc/par).
func (c *Controller) OidcPushedAuthorizationRequest(e echo.Context) error {
	req := e.Request()
	ctx := req.Context()

	ar, err := c.oauth2Provider.NewPushedAuthorizeRequest(ctx, req)
	if err != nil {
		return resterr.NewFositeError(resterr.FositePARError, e, c.oauth2Provider, err).WithAuthorizeRequester(ar)
	}

	var par PushedAuthorizationRequest

	if err = e.Bind(&par); err != nil {
		return resterr.NewFositeError(resterr.FositePARError, e, c.oauth2Provider, err).WithAuthorizeRequester(ar)
	}

	var ad common.AuthorizationDetails

	if err = json.Unmarshal([]byte(par.AuthorizationDetails), &ad); err != nil {
		return resterr.NewValidationError(resterr.InvalidValue, "authorization_details", err)
	}

	authorizationDetails, err := common.ValidateAuthorizationDetails(&ad)
	if err != nil {
		return err
	}

	r, err := c.issuerInteractionClient.PushAuthorizationDetails(ctx,
		issuer.PushAuthorizationDetailsJSONRequestBody{
			AuthorizationDetails: common.AuthorizationDetails{
				CredentialType: authorizationDetails.CredentialType,
				Format:         lo.ToPtr(string(authorizationDetails.Format)),
				Locations:      lo.ToPtr(authorizationDetails.Locations),
				Type:           authorizationDetails.Type,
			},
			OpState: par.OpState,
		},
	)
	if err != nil {
		return fmt.Errorf("push authorization details: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("push authorization details: status code %d", r.StatusCode)
	}

	resp, err := c.oauth2Provider.NewPushedAuthorizeResponse(ctx, ar, new(fosite.DefaultSession))
	if err != nil {
		return resterr.NewFositeError(resterr.FositePARError, e, c.oauth2Provider, err).WithAuthorizeRequester(ar)
	}

	c.oauth2Provider.WritePushedAuthorizeResponse(ctx, e.Response().Writer, ar, resp)

	return nil
}

// OidcAuthorize handles OIDC authorization request (GET /oidc/authorize).
func (c *Controller) OidcAuthorize(e echo.Context, params OidcAuthorizeParams) error { //nolint:funlen
	req := e.Request()
	ctx := req.Context()

	ar, err := c.oauth2Provider.NewAuthorizeRequest(ctx, req)
	if err != nil {
		return resterr.NewFositeError(resterr.FositeAuthorizeError, e, c.oauth2Provider, err).WithAuthorizeRequester(ar)
	}

	ses := &fosite.DefaultSession{
		Extra: map[string]interface{}{
			sessionOpStateKey:       params.OpState,
			authorizationDetailsKey: lo.FromPtr(params.AuthorizationDetails),
		},
	}

	var (
		credentialType string
		vcFormat       *string
	)

	scope := []string(ar.GetRequestedScopes())

	if params.AuthorizationDetails != nil {
		var authorizationDetails common.AuthorizationDetails

		if err = json.Unmarshal([]byte(*params.AuthorizationDetails), &authorizationDetails); err != nil {
			return resterr.NewValidationError(resterr.InvalidValue, "authorization_details", err)
		}

		if _, err = common.ValidateAuthorizationDetails(&authorizationDetails); err != nil {
			return err
		}

		credentialType = authorizationDetails.CredentialType
		vcFormat = authorizationDetails.Format
	} else {
		// using scope parameter to request credential type
		// https://openid.net/specs/openid-4-verifiable-credential-issuance-1_0.html#name-using-scope-parameter-to-re
		credentialType = scope[0]
	}

	r, err := c.issuerInteractionClient.PrepareAuthorizationRequest(ctx,
		issuer.PrepareAuthorizationRequestJSONRequestBody{
			AuthorizationDetails: &common.AuthorizationDetails{
				Type:           "openid_credential",
				CredentialType: credentialType,
				Format:         vcFormat,
			},
			OpState:      params.OpState,
			ResponseType: params.ResponseType,
			Scope:        lo.ToPtr(scope),
		},
	)
	if err != nil {
		return fmt.Errorf("prepare claim data authorization: %w", err)
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return fmt.Errorf("prepare claim data authorization: status code %d, %w",
			r.StatusCode,
			parseInteractionError(r.Body),
		)
	}

	var claimDataAuth issuer.PrepareClaimDataAuthorizationResponse

	if err = json.NewDecoder(r.Body).Decode(&claimDataAuth); err != nil {
		return fmt.Errorf("decode claim data authorization response: %w", err)
	}

	oauthConfig := oauth2.Config{
		ClientID:     claimDataAuth.AuthorizationRequest.ClientId,
		ClientSecret: claimDataAuth.AuthorizationRequest.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:   claimDataAuth.AuthorizationEndpoint,
			AuthStyle: oauth2.AuthStyleAutoDetect,
		},
		RedirectURL: c.issuerVCSPublicHost + "/oidc/redirect",
		Scopes:      claimDataAuth.AuthorizationRequest.Scope,
	}

	ar.(*fosite.AuthorizeRequest).State = params.OpState

	resp, err := c.oauth2Provider.NewAuthorizeResponse(ctx, ar, ses)
	if err != nil {
		return resterr.NewFositeError(resterr.FositeAuthorizeError, e, c.oauth2Provider, err).WithAuthorizeRequester(ar)
	}

	if err = c.stateStore.SaveAuthorizeState(
		ctx,
		params.OpState,
		&oidc4cistatestore.AuthorizeState{
			RedirectURI: ar.GetRedirectURI(),
			RespondMode: string(ar.GetResponseMode()),
			Header:      resp.GetHeader(),
			Parameters:  resp.GetParameters(),
		}); err != nil {
		return fmt.Errorf("save authorize state: %w", err)
	}

	var authCodeURL string
	if len(lo.FromPtr(claimDataAuth.PushedAuthorizationRequestEndpoint)) > 0 {
		authCodeURL, err = c.oAuth2Client.AuthCodeURLWithPAR(
			ctx,
			oauthConfig,
			*claimDataAuth.PushedAuthorizationRequestEndpoint,
			params.OpState,
			c.defaultHTTPClient,
		)
		if err != nil {
			return err
		}
	} else {
		authCodeURL = c.oAuth2Client.AuthCodeURL(ctx, oauthConfig, params.OpState)
	}

	return e.Redirect(http.StatusSeeOther, authCodeURL)
}

// OidcRedirect handles OIDC redirect (GET /oidc/redirect).
func (c *Controller) OidcRedirect(e echo.Context, params OidcRedirectParams) error {
	req := e.Request()
	ctx := req.Context()

	resp, err := c.stateStore.GetAuthorizeState(ctx, params.State)
	if err != nil {
		return apiUtil.WriteOutput(e)(nil, err)
	}

	storeResp, storeErr := c.issuerInteractionClient.StoreAuthorizationCodeRequest(ctx,
		issuer.StoreAuthorizationCodeRequestJSONRequestBody{
			Code:    params.Code,
			OpState: params.State,
		})
	if storeErr != nil {
		return storeErr
	}

	defer storeResp.Body.Close()

	if storeResp.StatusCode != http.StatusOK {
		return fmt.Errorf("store authorization code request: status code %d, %w",
			storeResp.StatusCode,
			parseInteractionError(storeResp.Body),
		)
	}

	responder := &fosite.AuthorizeResponse{}
	responder.Header = resp.Header
	responder.Parameters = resp.Parameters

	c.oauth2Provider.WriteAuthorizeResponse(ctx, e.Response().Writer, &fosite.AuthorizeRequest{
		RedirectURI:         resp.RedirectURI,
		ResponseMode:        fosite.ResponseModeType(resp.RespondMode),
		DefaultResponseMode: fosite.ResponseModeType(resp.RespondMode),
		State:               params.State,
	}, responder)

	return nil
}

// OidcToken handles OIDC token request (POST /oidc/token).
func (c *Controller) OidcToken(e echo.Context) error {
	req := e.Request()
	ctx := req.Context()

	ar, err := c.oauth2Provider.NewAccessRequest(ctx, req, new(fosite.DefaultSession))
	if err != nil {
		return resterr.NewFositeError(resterr.FositeAccessError, e, c.oauth2Provider, err).WithAccessRequester(ar)
	}

	session := ar.GetSession().(*fosite.DefaultSession) //nolint:errcheck
	if session.Extra == nil {
		session.Extra = make(map[string]interface{})
	}

	nonce := mustGenerateNonce()
	var txID string

	isPreAuthFlow := strings.EqualFold(e.FormValue("grant_type"), preAuthorizedCodeGrantType)
	if isPreAuthFlow { //nolint:nestif
		resp, preAuthorizeErr := c.oidcPreAuthorizedCode(
			ctx,
			e.FormValue("pre-authorized_code"),
			e.FormValue("user_pin"),
		)

		if preAuthorizeErr != nil {
			return preAuthorizeErr
		}

		txID = resp.TxId
	} else {
		exchangeResp, errExchange := c.issuerInteractionClient.ExchangeAuthorizationCodeRequest(
			ctx,
			issuer.ExchangeAuthorizationCodeRequestJSONRequestBody{
				OpState: ar.GetSession().(*fosite.DefaultSession).Extra[sessionOpStateKey].(string),
			},
		)
		if errExchange != nil {
			return fmt.Errorf("exchange authorization code request: %w", errExchange)
		}

		defer exchangeResp.Body.Close()

		if exchangeResp.StatusCode != http.StatusOK {
			return fmt.Errorf("exchange authorization code request: status code %d, %w",
				exchangeResp.StatusCode,
				parseInteractionError(exchangeResp.Body),
			)
		}

		var exchangeResult issuer.ExchangeAuthorizationCodeResponse

		if err = json.NewDecoder(exchangeResp.Body).Decode(&exchangeResult); err != nil {
			return fmt.Errorf("read exchange auth code response: %w", err)
		}
		txID = exchangeResult.TxId
	}

	c.setCNonceSession(session, nonce, txID)

	responder, err := c.oauth2Provider.NewAccessResponse(ctx, ar)
	if err != nil {
		return resterr.NewFositeError(resterr.FositeAccessError, e, c.oauth2Provider, err).WithAccessRequester(ar)
	}

	c.setCNonce(responder, nonce)

	c.oauth2Provider.WriteAccessResponse(ctx, e.Response().Writer, ar, responder)
	return nil
}

func (c *Controller) setCNonce(
	responder fosite.AccessResponder,
	nonce string,
) {
	responder.SetExtra("c_nonce", nonce)
	responder.SetExtra("c_nonce_expires_in", cNonceTTL.Seconds())
}

func (c *Controller) setCNonceSession(
	session *fosite.DefaultSession,
	nonce string,
	txID string,
) {
	session.Extra[txIDKey] = txID
	session.Extra[cNonceKey] = nonce
	session.Extra[cNonceExpiresAtKey] = time.Now().Add(cNonceTTL).Unix()
}

func mustGenerateNonce() string {
	b := make([]byte, cNonceSize)

	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}

	return base64.URLEncoding.EncodeToString(b)
}

// OidcCredential handles OIDC credential request (POST /oidc/credential).
func (c *Controller) OidcCredential(e echo.Context) error {
	req := e.Request()
	ctx := req.Context()

	var credentialRequest CredentialRequest

	if err := validateCredentialRequest(e, &credentialRequest); err != nil {
		return err
	}

	token := fosite.AccessTokenFromRequest(req)
	if token == "" {
		return resterr.NewUnauthorizedError(errors.New("missing access token"))
	}

	_, ar, err := c.oauth2Provider.IntrospectToken(ctx, token, fosite.AccessToken, new(fosite.DefaultSession))
	if err != nil {
		return resterr.NewUnauthorizedError(fmt.Errorf("introspect token: %w", err))
	}

	session := ar.GetSession().(*fosite.DefaultSession) //nolint:errcheck

	if err = validateProofClaims(credentialRequest.Proof.Jwt, session, c.jwtVerifier); err != nil {
		return resterr.NewValidationError(resterr.InvalidValue, "proof", err)
	}

	txID := session.Extra[txIDKey].(string) //nolint:errcheck

	resp, err := c.issuerInteractionClient.PrepareCredential(ctx,
		issuer.PrepareCredentialJSONRequestBody{
			TxId:   txID,
			Did:    lo.ToPtr(credentialRequest.Did),
			Type:   credentialRequest.Type,
			Format: credentialRequest.Format,
		},
	)
	if err != nil {
		return fmt.Errorf("prepare credential: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("prepare credential: status code %d, %w",
			resp.StatusCode,
			parseInteractionError(resp.Body),
		)
	}

	var result issuer.PrepareCredentialResult

	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode prepare credential result: %w", err)
	}

	nonce := mustGenerateNonce()

	session.Extra[cNonceKey] = nonce
	session.Extra[cNonceExpiresAtKey] = time.Now().Add(cNonceTTL).Unix()

	return apiUtil.WriteOutput(e)(CredentialResponse{
		Credential:      result.Credential,
		Format:          result.Format,
		CNonce:          lo.ToPtr(nonce),
		CNonceExpiresIn: lo.ToPtr(int(cNonceTTL.Seconds())),
	}, nil)
}

func validateCredentialRequest(e echo.Context, req *CredentialRequest) error {
	if err := e.Bind(req); err != nil {
		return err
	}

	_, err := common.ValidateVCFormat(common.VCFormat(lo.FromPtr(req.Format)))
	if err != nil {
		return resterr.NewValidationError(resterr.InvalidValue, "format", err)
	}

	if req.Proof == nil {
		return resterr.NewValidationError(resterr.InvalidValue, "proof", errors.New("missing proof type"))
	}

	if req.Proof.ProofType != "jwt" || req.Proof.Jwt == "" {
		return resterr.NewValidationError(resterr.InvalidValue, "proof", errors.New("invalid proof type"))
	}

	return nil
}

func validateProofClaims(rawJwt string, session *fosite.DefaultSession, verifier jose.SignatureVerifier) error {
	jws, err := jwt.Parse(rawJwt, jwt.WithSignatureVerifier(verifier))
	if err != nil {
		return fmt.Errorf("parse jwt: %w", err)
	}

	var claims JWTProofClaims

	if err = jws.DecodeClaims(&claims); err != nil {
		return fmt.Errorf("decode claims: %w", err)
	}

	if nonceExp := session.Extra[cNonceExpiresAtKey].(int64); nonceExp < time.Now().Unix() { //nolint:errcheck
		return errors.New("nonce expired")
	}

	if nonce := session.Extra[cNonceKey].(string); claims.Nonce != nonce { //nolint:errcheck
		return errors.New("invalid nonce")
	}

	return nil
}

// oidcPreAuthorizedCode handles pre-authorized code token request.
func (c *Controller) oidcPreAuthorizedCode(
	ctx context.Context,
	preAuthorizedCode string,
	userPin string,
) (*issuer.ValidatePreAuthorizedCodeResponse, error) {
	resp, err := c.issuerInteractionClient.ValidatePreAuthorizedCodeRequest(ctx,
		issuer.ValidatePreAuthorizedCodeRequestJSONRequestBody{
			PreAuthorizedCode: preAuthorizedCode,
			UserPin:           lo.ToPtr(userPin),
		})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validate pre-authorized code request: status code %d, %w",
			resp.StatusCode,
			parseInteractionError(resp.Body),
		)
	}

	var validateResponse issuer.ValidatePreAuthorizedCodeResponse
	if err = json.NewDecoder(resp.Body).Decode(&validateResponse); err != nil {
		return nil, err
	}

	return &validateResponse, nil
}

type interactionError struct {
	Code           string `json:"code"`
	Component      string `json:"component,omitempty"`
	Operation      string `json:"operation,omitempty"`
	IncorrectValue string `json:"incorrectValue,omitempty"`
	Message        string `json:"message,omitempty"`
}

func (e *interactionError) Error() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("code: %s", e.Code))

	if e.Component != "" {
		b.WriteString(fmt.Sprintf("; component: %s", e.Component))
	}

	if e.Operation != "" {
		b.WriteString(fmt.Sprintf("; operation: %s", e.Operation))
	}

	if e.IncorrectValue != "" {
		b.WriteString(fmt.Sprintf("; incorrect value: %s", e.IncorrectValue))
	}

	if e.Message != "" {
		b.WriteString(fmt.Sprintf("; message: %s", e.Message))
	}

	return b.String()
}

func parseInteractionError(reader io.Reader) error {
	b, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read body: %w", err)
	}

	var e interactionError

	if err = json.Unmarshal(b, &e); err != nil {
		return errors.New(string(b))
	}

	return &e
}
