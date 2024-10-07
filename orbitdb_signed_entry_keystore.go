package weshnet

import (
	"context"
	"encoding/hex"
	"sync"

	"github.com/libp2p/go-libp2p/core/crypto"

	"berty.tech/go-ipfs-log/keystore"
	"berty.tech/weshnet/v2/pkg/errcode"
)

type BertySignedKeyStore struct {
	sync.Map
}

func (s *BertySignedKeyStore) SetKey(pk crypto.PrivKey) error {
	pubKeyBytes, err := pk.GetPublic().Raw()
	if err != nil {
		return errcode.ErrCode_TODO.Wrap(err)
	}

	keyID := hex.EncodeToString(pubKeyBytes)

	s.Store(keyID, pk)

	return nil
}

func (s *BertySignedKeyStore) HasKey(_ context.Context, id string) (bool, error) {
	_, ok := s.Load(id)

	return ok, nil
}

func (s *BertySignedKeyStore) CreateKey(ctx context.Context, id string) (crypto.PrivKey, error) {
	return s.GetKey(ctx, id)
}

func (s *BertySignedKeyStore) GetKey(_ context.Context, id string) (crypto.PrivKey, error) {
	if privKey, ok := s.Load(id); ok {
		if pk, ok := privKey.(crypto.PrivKey); ok {
			return pk, nil
		}
	}

	return nil, errcode.ErrCode_ErrGroupMemberUnknownGroupID
}

func (s *BertySignedKeyStore) Sign(privKey crypto.PrivKey, bytes []byte) ([]byte, error) {
	return privKey.Sign(bytes)
}

func (s *BertySignedKeyStore) Verify(signature []byte, publicKey crypto.PubKey, data []byte) error {
	ok, err := publicKey.Verify(data, signature)
	if err != nil {
		return err
	}

	if !ok {
		return errcode.ErrCode_ErrGroupMemberLogEventSignature
	}

	return nil
}

func (s *BertySignedKeyStore) getIdentityProvider() *bertySignedIdentityProvider {
	return &bertySignedIdentityProvider{
		keyStore: s,
	}
}

var _ keystore.Interface = (*BertySignedKeyStore)(nil)
