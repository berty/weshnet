package weshnet

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"berty.tech/weshnet/pkg/bertyvcissuer"
	"berty.tech/weshnet/pkg/cryptoutil"
	"berty.tech/weshnet/pkg/errcode"
	"berty.tech/weshnet/pkg/protocoltypes"
)

func (s *service) CredentialVerificationServiceInitFlow(ctx context.Context, request *protocoltypes.CredentialVerificationServiceInitFlow_Request) (*protocoltypes.CredentialVerificationServiceInitFlow_Reply, error) {
	s.lock.Lock()
	s.vcClient = bertyvcissuer.NewClient(request.ServiceUrl)
	client := s.vcClient
	s.lock.Unlock()

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	// TODO: allow selection of alt-scoped keys
	// TODO: avoid exporting account keys
	pkRaw, err := s.accountGroupCtx.ownMemberDevice.Member().Raw()
	if err != nil {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	if !bytes.Equal(pkRaw, request.PublicKey) {
		return nil, errcode.ErrCode_ErrInvalidInput
	}

	url, err := client.Init(ctx, request.Link, cryptoutil.NewFuncSigner(s.accountGroupCtx.ownMemberDevice.Member(), s.accountGroupCtx.ownMemberDevice.MemberSign))
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return &protocoltypes.CredentialVerificationServiceInitFlow_Reply{
		Url:       url,
		SecureUrl: strings.HasPrefix(url, "https://"),
	}, nil
}

func (s *service) CredentialVerificationServiceCompleteFlow(ctx context.Context, request *protocoltypes.CredentialVerificationServiceCompleteFlow_Request) (*protocoltypes.CredentialVerificationServiceCompleteFlow_Reply, error) {
	s.lock.Lock()
	client := s.vcClient
	s.lock.Unlock()

	if client == nil {
		return nil, errcode.ErrCode_ErrInvalidInput.Wrap(fmt.Errorf("a verification flow needs to be started first"))
	}

	credentials, identifier, parsedCredential, err := client.Complete(request.CallbackUri)
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	_, err = s.accountGroupCtx.metadataStore.SendAccountVerifiedCredentialAdded(ctx, &protocoltypes.AccountVerifiedCredentialRegistered{
		VerifiedCredential: credentials,
		RegistrationDate:   parsedCredential.Issued.UnixNano(),
		ExpirationDate:     parsedCredential.Expired.UnixNano(),
		Identifier:         identifier,
		Issuer:             parsedCredential.Issuer.ID,
	})
	if err != nil {
		return nil, errcode.ErrCode_ErrInternal.Wrap(err)
	}

	return &protocoltypes.CredentialVerificationServiceCompleteFlow_Reply{
		Identifier: identifier,
	}, nil
}

func (s *service) VerifiedCredentialsList(request *protocoltypes.VerifiedCredentialsList_Request, server protocoltypes.ProtocolService_VerifiedCredentialsListServer) error {
	now := time.Now().UnixNano()
	credentials := s.accountGroupCtx.metadataStore.ListVerifiedCredentials()

	for _, credential := range credentials {
		if request.FilterIdentifier != "" && credential.Identifier != request.FilterIdentifier {
			continue
		}

		if request.ExcludeExpired && credential.ExpirationDate < now {
			continue
		}

		if request.FilterIssuer != "" && credential.Issuer != request.FilterIssuer {
			continue
		}

		if err := server.Send(&protocoltypes.VerifiedCredentialsList_Reply{
			Credential: credential,
		}); err != nil {
			return errcode.ErrCode_ErrStreamWrite.Wrap(err)
		}
	}

	return nil
}
