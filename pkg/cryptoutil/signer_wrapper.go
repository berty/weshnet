package cryptoutil

import (
	stdcrypto "crypto"
	"io"

	libp2p_ci "github.com/libp2p/go-libp2p/core/crypto"
)

func NewFuncSigner(key libp2p_ci.PubKey, signer func([]byte) ([]byte, error)) stdcrypto.Signer {
	return &funcSigner{
		pubKey: key,
		signer: signer,
	}
}

type funcSigner struct {
	pubKey libp2p_ci.PubKey
	signer func([]byte) ([]byte, error)
}

func (f *funcSigner) Public() stdcrypto.PublicKey {
	return f.pubKey
}

func (f *funcSigner) Sign(_ io.Reader, digest []byte, _ stdcrypto.SignerOpts) (signature []byte, err error) {
	return f.signer(digest)
}
