package bertyprotocol

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"berty.tech/berty/v2/go/internal/cryptoutil"
	"berty.tech/berty/v2/go/pkg/bertytypes"
	"berty.tech/berty/v2/go/pkg/errcode"
)

const AuthResponseType = "code"
const AuthGrantType = "authorization_code"
const AuthRedirect = "berty://services-auth/"
const AuthClientID = "berty"
const AuthCodeChallengeMethod = "S256"

type authExchangeResponse struct {
	AccessToken      string            `json:"access_token"`
	Scope            string            `json:"scope"`
	Error            string            `json:"error"`
	ErrorDescription string            `json:"error_description"`
	Services         map[string]string `json:"services"`
}

type authSession struct {
	state        string
	codeVerifier string
	baseURL      string
}

func newAuthSession(baseURL string) (*authSession, error) {
	state, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, err
	}

	codeVerifier, err := cryptoutil.GenerateNonce()
	if err != nil {
		return nil, err
	}

	codeVerifierBytes := make([]byte, cryptoutil.NonceSize)
	for i, c := range codeVerifier {
		codeVerifierBytes[i] = c
	}

	stateBytes := make([]byte, cryptoutil.NonceSize)
	for i, c := range state {
		stateBytes[i] = c
	}

	auth := &authSession{
		baseURL:      baseURL,
		state:        base64.RawURLEncoding.EncodeToString(stateBytes),
		codeVerifier: base64.RawURLEncoding.EncodeToString(codeVerifierBytes),
	}

	return auth, nil
}

func (s *service) authInitURL(baseURL string) (string, error) {
	if !strings.HasSuffix(baseURL, "/") {
		baseURL += "/"
	}

	auth, err := newAuthSession(baseURL)
	if err != nil {
		return "", err
	}

	s.authSession.Store(auth)

	codeChallengeArr := sha256.Sum256([]byte(auth.codeVerifier))
	codeChallenge := make([]byte, sha256.Size)
	for i, c := range codeChallengeArr {
		codeChallenge[i] = c
	}

	return fmt.Sprintf("%sauthorize?response_type=%s&client_id=%s&redirect_uri=%s&state=%s&code_challenge=%s&code_challenge_method=%s",
		baseURL,
		AuthResponseType,
		AuthClientID,
		url.QueryEscape(AuthRedirect),
		auth.state,
		base64.RawURLEncoding.EncodeToString(codeChallenge),
		AuthCodeChallengeMethod,
	), nil
}

func (s *service) AuthServiceCompleteFlow(ctx context.Context, request *bertytypes.AuthServiceCompleteFlow_Request) (*bertytypes.AuthServiceCompleteFlow_Reply, error) {
	u, err := url.Parse(request.CallbackURL)
	if err != nil {
		return nil, err
	}

	if e := u.Query().Get("error"); e != "" {
		return nil, errcode.ErrServicesAuthServer.Wrap(fmt.Errorf("got error: %s (%s)", e, u.Query().Get("error_description")))
	}

	code, state := u.Query().Get("code"), u.Query().Get("state")

	authUntyped := s.authSession.Load()
	if authUntyped == nil {
		return nil, errcode.ErrServicesAuthNotInitialized
	}

	auth, ok := authUntyped.(*authSession)
	if !ok {
		return nil, errcode.ErrServicesAuthNotInitialized
	}

	if auth.state != state {
		return nil, errcode.ErrServicesAuthWrongState
	}

	res, err := http.PostForm(fmt.Sprintf("%soauth/token", auth.baseURL), url.Values{
		"grant_type":    {AuthGrantType},
		"code":          {code},
		"client_id":     {AuthClientID},
		"code_verifier": {auth.codeVerifier},
	})

	if err != nil {
		return nil, errcode.ErrStreamWrite.Wrap(err)
	}

	defer res.Body.Close()

	if res.StatusCode >= 300 {
		return nil, errcode.ErrServicesAuthInvalidResponse.Wrap(fmt.Errorf("invalid status code %d", res.StatusCode))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errcode.ErrStreamRead.Wrap(err)
	}

	resMsg := &authExchangeResponse{}
	if err := json.Unmarshal(body, &resMsg); err != nil {
		return nil, errcode.ErrDeserialization.Wrap(err)
	}

	if resMsg.Error != "" {
		return nil, errcode.ErrServicesAuthServer.Wrap(err)
	}

	if resMsg.AccessToken == "" {
		return nil, errcode.ErrServicesAuthInvalidResponse.Wrap(fmt.Errorf("missing access token in response"))
	}

	if len(resMsg.Services) == 0 {
		return nil, errcode.ErrServicesAuthInvalidResponse.Wrap(fmt.Errorf("no services returned along token"))
	}

	services := make([]*bertytypes.ServiceTokenSupportedService, len(resMsg.Services))
	i := 0
	for k, v := range resMsg.Services {
		services[i] = &bertytypes.ServiceTokenSupportedService{
			ServiceType:     k,
			ServiceEndpoint: v,
		}
		i++
	}

	if _, err := s.accountGroup.metadataStore.SendAccountServiceTokenAdded(ctx, &bertytypes.ServiceToken{
		Token:             resMsg.AccessToken,
		AuthenticationURL: auth.baseURL,
		SupportedServices: services,
		Expiration:        -1,
	}); err != nil {
		return nil, err
	}

	return &bertytypes.AuthServiceCompleteFlow_Reply{}, nil
}

func (s *service) AuthServiceInitFlow(ctx context.Context, request *bertytypes.AuthServiceInitFlow_Request) (*bertytypes.AuthServiceInitFlow_Reply, error) {
	parsedAuthURL, err := url.Parse(request.AuthURL)
	if err != nil {
		return nil, errcode.ErrServicesAuthInvalidURL
	}

	secure := true

	if parsedAuthURL.Scheme == "http" {
		secure = false
	} else if parsedAuthURL.Scheme != "https" {
		return nil, errcode.ErrServicesAuthInvalidURL
	}

	u, err := s.authInitURL(request.AuthURL)
	if err != nil {
		return nil, err
	}

	return &bertytypes.AuthServiceInitFlow_Reply{
		URL:       u,
		SecureURL: secure,
	}, nil
}

func (s *service) ServicesTokenList(request *bertytypes.ServicesTokenList_Request, server ProtocolService_ServicesTokenListServer) error {
	for _, t := range s.accountGroup.metadataStore.listServiceTokens() {
		if server.Context().Err() != nil {
			break
		}

		if err := server.Send(&bertytypes.ServicesTokenList_Reply{
			TokenID: t.TokenID(),
			Service: t,
		}); err != nil {
			return err
		}
	}

	return nil
}
