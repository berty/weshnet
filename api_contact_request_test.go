package weshnet

import (
	"context"
	"testing"
	"time"

	libp2p_mocknet "github.com/berty/go-libp2p-mock"
	"github.com/stretchr/testify/require"

	"berty.tech/weshnet/pkg/protocoltypes"
	"berty.tech/weshnet/pkg/testutil"
)

func TestShareContact(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	logger, cleanup := testutil.Logger(t)
	defer cleanup()

	opts := TestingOpts{
		Mocknet: libp2p_mocknet.New(),
		Logger:  logger,
	}

	pts, cleanup := NewTestingProtocolWithMockedPeers(ctx, t, &opts, nil, 2)
	defer cleanup()

	binaryContact, err := pts[0].Client.ShareContact(ctx, &protocoltypes.ShareContact_Request{})
	require.NoError(t, err)

	// Check that ShareContact reset the contact request reference and enabled contact requests.
	contactRequestRef, err := pts[0].Client.ContactRequestReference(ctx,
		&protocoltypes.ContactRequestReference_Request{})
	require.NoError(t, err)

	require.NotEqual(t, 0, len(contactRequestRef.PublicRendezvousSeed))
	require.Equal(t, true, contactRequestRef.Enabled)

	// Decode.
	contact, err := pts[0].Client.DecodeContact(ctx, &protocoltypes.DecodeContact_Request{
		EncodedContact: binaryContact.EncodedContact,
	})
	require.NoError(t, err)

	// Check for the expected info.
	config, err := pts[0].Client.ServiceGetConfiguration(ctx,
		&protocoltypes.ServiceGetConfiguration_Request{})
	require.NoError(t, err)
	require.Equal(t, contact.Contact.PK, config.AccountPK)
	require.Equal(t, contact.Contact.PublicRendezvousSeed, contactRequestRef.PublicRendezvousSeed)
}
