package secretstore_test

import (
	crand "crypto/rand"
	"testing"

	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/stretchr/testify/assert"

	"berty.tech/weshnet/v2/pkg/protocoltypes"
	"berty.tech/weshnet/v2/pkg/secretstore"
)

func Test_New_AccountPrivKey_AccountProofPrivKey(t *testing.T) {
	acc, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc)

	sk1, skProof1, err := acc.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk1)
	assert.NotNil(t, skProof1)

	sk2, skProof2, err := acc.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk2)
	assert.NotNil(t, skProof2)

	assert.Equal(t, sk1, sk2)
	assert.Equal(t, skProof1, skProof2)

	assert.NotEqual(t, sk1, skProof1)
	assert.NotEqual(t, sk1, skProof2)
	assert.NotEqual(t, sk2, skProof1)
	assert.NotEqual(t, sk2, skProof2)
}

func Test_ExportAccountKeys_ImportAccountKeys(t *testing.T) {
	acc1, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc1)

	sk1, skProof1, err := acc1.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk1)
	assert.NotNil(t, skProof1)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc2)

	// Testing with a nil value
	{
		assert.Error(t, acc2.ImportAccountKeys(nil, skProof1))
		assert.Error(t, acc2.ImportAccountKeys(sk1, nil))
	}

	// Testing with an unsupported key format
	{
		invalidPriv, _, err := crypto.GenerateSecp256k1Key(crand.Reader)
		assert.NoError(t, err)

		invalidPrivBytes, err := crypto.MarshalPrivateKey(invalidPriv)
		assert.NoError(t, err)

		assert.Error(t, acc2.ImportAccountKeys(sk1, invalidPrivBytes))
		assert.Error(t, acc2.ImportAccountKeys(invalidPrivBytes, skProof1))
	}

	// Testing with an invalid key format
	{
		garbageBytes := []byte("garbage")
		assert.Error(t, acc2.ImportAccountKeys(sk1, garbageBytes))
		assert.Error(t, acc2.ImportAccountKeys(garbageBytes, skProof1))
	}

	// Testing with account and proof key being the same
	{
		assert.Error(t, acc2.ImportAccountKeys(sk1, sk1))
	}

	// Valid test case
	{
		assert.NoError(t, acc2.ImportAccountKeys(sk1, skProof1))
	}

	// Attempting to import keys again
	{
		assert.Error(t, acc2.ImportAccountKeys(sk1, skProof1))
	}

	sk2, skProof2, err := acc1.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk2)
	assert.NotNil(t, skProof2)

	assert.Equal(t, sk1, sk2)
	assert.Equal(t, skProof1, skProof2)
	assert.NotEqual(t, sk1, skProof1)
	assert.NotEqual(t, sk1, skProof2)
	assert.NotEqual(t, sk2, skProof1)
	assert.NotEqual(t, sk2, skProof2)
}

func Test_DevicePrivKey(t *testing.T) {
	acc1, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc1)

	sk1, skProof1, err := acc1.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk1)
	assert.NotNil(t, skProof1)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc1)

	err = acc2.ImportAccountKeys(sk1, skProof1)
	assert.NoError(t, err)

	sk2, skProof2, err := acc2.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk2)
	assert.NotNil(t, skProof2)

	accGroup1, memberDevice1, err := acc1.GetGroupForAccount()
	assert.NoError(t, err)
	assert.NotNil(t, accGroup1)

	memberDevice2, err := acc2.GetOwnMemberDeviceForGroup(accGroup1)
	assert.NoError(t, err)
	assert.NotNil(t, memberDevice2)

	assert.True(t, memberDevice1.Member().Equals(memberDevice2.Member()))
	assert.False(t, memberDevice1.Device().Equals(memberDevice2.Device()))
}

func Test_ContactGroupPrivKey(t *testing.T) {
	acc1, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc1)

	_, acc1MemberDevice, err := acc1.GetGroupForAccount()
	assert.NoError(t, err)
	assert.NotNil(t, acc1MemberDevice)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc2)

	_, acc2MemberDevice, err := acc2.GetGroupForAccount()
	assert.NoError(t, err)
	assert.NotNil(t, acc2MemberDevice)

	grp1, err := acc1.GetGroupForContact(acc2MemberDevice.Member())
	assert.NoError(t, err)
	assert.NotNil(t, grp1)

	grp2, err := acc2.GetGroupForContact(acc1MemberDevice.Member())
	assert.NoError(t, err)
	assert.NotNil(t, grp2)

	assert.Equal(t, grp1.PublicKey, grp2.PublicKey)
	assert.Equal(t, grp1.Secret, grp2.Secret)
	assert.Equal(t, grp1.GroupType, grp2.GroupType)
}

func Test_MemberDeviceForGroup_multimember(t *testing.T) {
	acc1, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc1)

	sk1, skProof1, err := acc1.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk1)
	assert.NotNil(t, skProof1)

	acc2, err := secretstore.NewInMemSecretStore(nil)
	assert.NoError(t, err)
	assert.NotNil(t, acc2)

	err = acc2.ImportAccountKeys(sk1, skProof1)
	assert.NoError(t, err)

	sk2, skProof2, err := acc2.ExportAccountKeysForBackup()
	assert.NoError(t, err)
	assert.NotNil(t, sk2)
	assert.NotNil(t, skProof2)

	g, _, err := protocoltypes.NewGroupMultiMember()
	assert.NoError(t, err)

	omd1, err := acc1.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	omd2, err := acc2.GetOwnMemberDeviceForGroup(g)
	assert.NoError(t, err)

	omd1M := omd1.Member()
	omd2M := omd2.Member()
	omd1D := omd1.Device()
	omd2D := omd2.Device()

	assert.True(t, omd1M.Equals(omd2M))
	assert.False(t, omd1D.Equals(omd2D))
}
