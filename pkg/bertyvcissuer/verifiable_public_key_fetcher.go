package bertyvcissuer

import (
	"crypto/ed25519"
	"fmt"
	"slices"
	"strings"

	"github.com/hyperledger/aries-framework-go/pkg/doc/signature/verifier"
	"github.com/hyperledger/aries-framework-go/pkg/kms"
	"github.com/multiformats/go-multibase"

	"berty.tech/weshnet/v2/pkg/errcode"
)

func embeddedPublicKeyFetcher(issuerID string, allowList []string) (*verifier.PublicKey, error) {
	if !strings.HasPrefix(issuerID, "did:key:z6Mk") {
		return nil, fmt.Errorf("unexpected key format")
	}

	if len(allowList) > 0 {
		found := slices.Contains(allowList, issuerID)

		if !found {
			return nil, errcode.ErrCode_ErrServicesDirectoryInvalidVerifiedCredentialID.Wrap(fmt.Errorf("issuer is not allowed"))
		}
	}

	_, rawData, err := multibase.Decode(issuerID[8:])
	if err != nil {
		return nil, err
	}

	if len(rawData) != ed25519.PublicKeySize+2 {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	return &verifier.PublicKey{
		Type:  kms.ED25519,
		Value: rawData[2:],
		JWK:   nil,
	}, nil
}

// nolint:revive
func EmbeddedPublicKeyFetcher(issuerID, keyID string) (*verifier.PublicKey, error) {
	return embeddedPublicKeyFetcher(issuerID, nil)
}

// nolint:revive
func EmbeddedPublicKeyFetcherAllowList(allowList []string) func(issuerID, keyID string) (*verifier.PublicKey, error) {
	return func(issuerID, keyID string) (*verifier.PublicKey, error) {
		return embeddedPublicKeyFetcher(issuerID, allowList)
	}
}
