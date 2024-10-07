package handshake

import (
	"context"
	"errors"

	p2pcrypto "github.com/libp2p/go-libp2p/core/crypto"
	"go.uber.org/zap"
	"golang.org/x/crypto/nacl/box"
	"google.golang.org/protobuf/proto"

	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protoio"
	"berty.tech/weshnet/pkg/tyber"
)

// ResponseUsingReaderWriter handle the handshake inited by the requester, using provided io reader and writer
func ResponseUsingReaderWriter(ctx context.Context, logger *zap.Logger, reader protoio.Reader, writer protoio.Writer, ownAccountID p2pcrypto.PrivKey) (p2pcrypto.PubKey, error) {
	hc := &handshakeContext{
		reader:          reader,
		writer:          writer,
		ownAccountID:    ownAccountID,
		sharedEphemeral: &[cryptoutil.KeySize]byte{},
	}

	// Handshake steps on responder side (see comments below)
	if err := hc.receiveRequesterHello(); err != nil {
		return nil, errcode.ErrCode_ErrHandshakeRequesterHello.Wrap(err)
	}
	tyber.LogStep(ctx, logger, "Received hello", hc.toTyberStepMutator(), tyber.ForceReopen)
	if err := hc.sendResponderHello(); err != nil {
		return nil, errcode.ErrCode_ErrHandshakeResponderHello.Wrap(err)
	}
	tyber.LogStep(ctx, logger, "Sent hello", hc.toTyberStepMutator(), tyber.ForceReopen)
	if err := hc.receiveRequesterAuthenticate(); err != nil {
		return nil, errcode.ErrCode_ErrHandshakeRequesterAuthenticate.Wrap(err)
	}
	tyber.LogStep(ctx, logger, "Received authenticate", hc.toTyberStepMutator(), tyber.ForceReopen)
	if err := hc.sendResponderAccept(); err != nil {
		return nil, errcode.ErrCode_ErrHandshakeResponderAccept.Wrap(err)
	}
	tyber.LogStep(ctx, logger, "Sent accept", hc.toTyberStepMutator(), tyber.ForceReopen)
	if err := hc.receiveRequesterAcknowledge(); err != nil {
		return nil, errcode.ErrCode_ErrHandshakeRequesterAcknowledge.Wrap(err)
	}
	tyber.LogStep(ctx, logger, "Received acknowledge", hc.toTyberStepMutator(), tyber.ForceReopen)

	return hc.peerAccountID, nil
}

// 1st step - Responder receives: a
func (hc *handshakeContext) receiveRequesterHello() error {
	if err := hc.receivePeerEphemeralPubKey(); err != nil {
		return errcode.ErrCode_ErrHandshakePeerEphemeralKeyRecv.Wrap(err)
	}

	return nil
}

// 2nd step - Responder sends: b
func (hc *handshakeContext) sendResponderHello() error {
	if err := hc.generateOwnEphemeralAndSendPubKey(); err != nil {
		return errcode.ErrCode_ErrHandshakeOwnEphemeralKeyGenSend.Wrap(err)
	}

	// Compute shared key from Ephemeral keys
	box.Precompute(hc.sharedEphemeral, hc.peerEphemeral, hc.ownEphemeral)

	return nil
}

// 3rd step - Responder receives: box[a.b|a.B](A,sig[A](a.b))
func (hc *handshakeContext) receiveRequesterAuthenticate() error {
	var (
		boxEnvelope BoxEnvelope
		request     RequesterAuthenticatePayload
	)

	// Receive BoxEnvelope from requester
	if err := hc.reader.ReadMsg(&boxEnvelope); err != nil {
		return errcode.ErrCode_ErrStreamRead.Wrap(err)
	}

	// Compute box key and open marshaled RequesterAuthenticatePayload using
	// constant nonce (see handshake.go)
	boxKey, err := hc.computeRequesterAuthenticateBoxKey(false)
	if err != nil {
		return errcode.ErrCode_ErrHandshakeRequesterAuthenticateBoxKeyGen.Wrap(err)
	}
	requestBytes, _ := box.OpenAfterPrecomputation(
		nil,
		boxEnvelope.Box,
		&nonceRequesterAuthenticate,
		boxKey,
	)
	if requestBytes == nil {
		err := errors.New("box opening failed")
		return errcode.ErrCode_ErrCryptoDecrypt.Wrap(err)
	}

	// Unmarshal RequesterAuthenticatePayload and RequesterAccountId
	err = proto.Unmarshal(requestBytes, &request)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}
	hc.peerAccountID, err = p2pcrypto.UnmarshalPublicKey(request.RequesterAccountId)
	if err != nil {
		return errcode.ErrCode_ErrDeserialization.Wrap(err)
	}

	// Verify proof (shared_a_b signed by peer's AccountID)
	valid, err := hc.peerAccountID.Verify(
		hc.sharedEphemeral[:],
		request.RequesterAccountSig,
	)
	if err != nil {
		return errcode.ErrCode_ErrCryptoSignatureVerification.Wrap(err)
	} else if !valid {
		return errcode.ErrCode_ErrCryptoSignatureVerification
	}

	return nil
}

// 4th step - Responder sends: box[a.b|A.B](sig[B](a.b))
func (hc *handshakeContext) sendResponderAccept() error {
	var (
		response ResponderAcceptPayload
		err      error
	)

	// Set proof (shared_a_b signed by own AccountID) in ResponderAcceptPayload
	// before marshaling it
	response.ResponderAccountSig, err = hc.ownAccountID.Sign(hc.sharedEphemeral[:])
	if err != nil {
		return errcode.ErrCode_ErrCryptoSignature.Wrap(err)
	}
	responseBytes, err := proto.Marshal(&response)
	if err != nil {
		return errcode.ErrCode_ErrSerialization.Wrap(err)
	}

	// Compute box key and seal marshaled ResponderAcceptPayload using
	// constant nonce (see handshake.go)
	boxKey, err := hc.computeResponderAcceptBoxKey()
	if err != nil {
		return errcode.ErrCode_ErrHandshakeResponderAcceptBoxKeyGen.Wrap(err)
	}
	boxContent := box.SealAfterPrecomputation(
		nil,
		responseBytes,
		&nonceResponderAccept,
		boxKey,
	)

	// Send BoxEnvelope to requester
	if err = hc.writer.WriteMsg(&BoxEnvelope{Box: boxContent}); err != nil {
		return errcode.ErrCode_ErrStreamWrite.Wrap(err)
	}

	return nil
}

// 5th step - Responder receives: ok
func (hc *handshakeContext) receiveRequesterAcknowledge() error {
	var acknowledge RequesterAcknowledgePayload

	// Receive Acknowledge from requester
	if err := hc.reader.ReadMsg(&acknowledge); err != nil {
		return errcode.ErrCode_ErrStreamRead.Wrap(err)
	}

	if !acknowledge.Success {
		return errcode.ErrCode_ErrInvalidInput
	}

	return nil
}
