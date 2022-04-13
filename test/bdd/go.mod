// Copyright SecureKey Technologies Inc. All Rights Reserved.
//
// SPDX-License-Identifier: Apache-2.0

module github.com/trustbloc/edge-service/test/bdd

replace github.com/trustbloc/edge-service => ../..

go 1.16

require (
	github.com/Microsoft/hcsshim v0.8.21 // indirect
	github.com/btcsuite/btcutil v1.0.3-0.20201208143702-a53e38424cce
	github.com/containerd/containerd v1.5.5 // indirect
	github.com/cucumber/godog v0.9.0
	github.com/docker/docker v20.10.8+incompatible // indirect
	github.com/fsouza/go-dockerclient v1.6.0
	github.com/go-openapi/runtime v0.19.26
	github.com/go-openapi/strfmt v0.20.0
	github.com/google/uuid v1.3.0
	github.com/hyperledger/aries-framework-go v0.1.9-0.20220412155017-81442062e607
	github.com/hyperledger/aries-framework-go-ext/component/vdr/orb v1.0.0-rc.1
	github.com/hyperledger/aries-framework-go/component/storageutil v0.0.0-20220330140627-07042d78580c
	github.com/hyperledger/aries-framework-go/spi v0.0.0-20220330140627-07042d78580c
	github.com/igor-pavlenko/httpsignatures-go v0.0.23
	github.com/moby/sys/mount v0.2.0 // indirect
	github.com/tidwall/gjson v1.6.7
	github.com/trustbloc/edge-core v0.1.8
	github.com/trustbloc/edge-service v0.0.0-00010101000000-000000000000
	github.com/trustbloc/edv v0.1.8
)
