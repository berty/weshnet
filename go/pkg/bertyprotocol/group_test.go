package bertyprotocol

import (
	"crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/stretchr/testify/require"
)

func TestGetGroupForContact(t *testing.T) {
	sk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	g, err := GetGroupForContact(sk)
	require.NoError(t, err)

	require.Equal(t, g.GroupType, GroupTypeContact)
	require.Equal(t, len(g.PublicKey), 32)
	require.Equal(t, len(g.Secret), 32)
}

func TestGetKeysForGroupOfContact(t *testing.T) {
	sk, _, err := crypto.GenerateEd25519Key(rand.Reader)
	require.NoError(t, err)

	sk1, sk2, err := getKeysForGroupOfContact(sk)
	require.NoError(t, err)

	require.NotNil(t, sk1)
	require.NotNil(t, sk2)
	require.False(t, sk1.Equals(sk2))
}
