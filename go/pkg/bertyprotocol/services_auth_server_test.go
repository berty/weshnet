package bertyprotocol

import (
	"bytes"
	"html/template"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func mustReadAllBytes(t *testing.T, reader io.ReadCloser) []byte {
	data, err := ioutil.ReadAll(reader)
	_ = reader.Close()
	require.NoError(t, err)

	return data
}

func TestNewAuthTokenServer(t *testing.T) {
	services := map[string]string{
		"service": "servicehost:1234",
	}
	secret, _, sk := helperGenerateTokenIssuerSecrets(t)

	ats, err := NewAuthTokenServer(secret, nil, services, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, ats)

	ats, err = NewAuthTokenServer(secret, sk, nil, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, ats)

	ats, err = NewAuthTokenServer(secret, sk, map[string]string{}, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, ats)

	ats, err = NewAuthTokenServer(nil, sk, services, zap.NewNop())
	require.Error(t, err)
	require.Nil(t, ats)

	ats, err = NewAuthTokenServer(secret, sk, services, nil)
	require.NoError(t, err)
	require.NotNil(t, ats)

	ats, err = NewAuthTokenServer(secret, sk, services, zap.NewNop())
	require.NoError(t, err)
	require.NotNil(t, ats)

	mux := ats.serveMux()
	server := httptest.NewServer(mux)

	defer server.Close()

	res, err := server.Client().Get(server.URL)
	require.NoError(t, err)
	require.Equal(t, 404, res.StatusCode)

	authorizeURL, err := url.Parse(server.URL)
	require.NoError(t, err)

	authorizeURL.Path = AuthHTTPPathAuthorize

	res, err = server.Client().Get(authorizeURL.String())
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, string(templateMustExec(t, templateAuthTokenServerRedirect, map[string]template.URL{
		"URL": "?error=invalid_request&error_description=unexpected+value+for+redirect_uri",
	})), string(mustReadAllBytes(t, res.Body)))

	responseRedirectURI, err := url.Parse(AuthRedirect)
	require.NoError(t, err)
	setURLParam(responseRedirectURI, "error", "invalid_request")
	setURLParam(responseRedirectURI, "error_description", "unexpected value for response_type")

	setURLParam(authorizeURL, "redirect_uri", AuthRedirect)

	testAuthorizeQueryURLAndCompareResponse(t, server, authorizeURL, responseRedirectURI)

	setURLParam(authorizeURL, "response_type", AuthResponseType)
	setURLParam(responseRedirectURI, "error_description", "unexpected value for client_id")

	testAuthorizeQueryURLAndCompareResponse(t, server, authorizeURL, responseRedirectURI)

	setURLParam(authorizeURL, "client_id", AuthClientID)
	setURLParam(responseRedirectURI, "error_description", "unexpected value for code_challenge_method")

	testAuthorizeQueryURLAndCompareResponse(t, server, authorizeURL, responseRedirectURI)

	setURLParam(authorizeURL, "code_challenge_method", AuthCodeChallengeMethod)
	setURLParam(responseRedirectURI, "error_description", "unexpected value for state")

	testAuthorizeQueryURLAndCompareResponse(t, server, authorizeURL, responseRedirectURI)

	setURLParam(authorizeURL, "state", "some_state")
	setURLParam(responseRedirectURI, "error_description", "unexpected value for code_challenge")

	testAuthorizeQueryURLAndCompareResponse(t, server, authorizeURL, responseRedirectURI)

	setURLParam(authorizeURL, "code_challenge", "some_code_challenge")
	setURLParam(responseRedirectURI, "error_description", "unexpected value for code_challenge")

	res, err = server.Client().Get(authorizeURL.String())
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, templateAuthTokenServerAuthorizeButton, string(mustReadAllBytes(t, res.Body)))

	res, err = server.Client().Post(authorizeURL.String(), "multipart/form-data", bytes.NewBuffer([]byte("")))
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Contains(t, string(mustReadAllBytes(t, res.Body)), "?code=eyJ")
}

func testAuthorizeQueryURLAndCompareResponse(t *testing.T, server *httptest.Server, queryURL *url.URL, redirectURL *url.URL) {
	testAuthorizeGetQueryURLAndCompareResponse(t, server, queryURL, redirectURL)
	testAuthorizePostQueryURLAndCompareResponse(t, server, queryURL, redirectURL)
}

func testAuthorizeGetQueryURLAndCompareResponse(t *testing.T, server *httptest.Server, queryURL *url.URL, redirectURL *url.URL) {
	res, err := server.Client().Get(queryURL.String())
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, string(templateMustExec(t, templateAuthTokenServerRedirect, map[string]template.URL{
		"URL": template.URL(redirectURL.String()),
	})), string(mustReadAllBytes(t, res.Body)))
}

func testAuthorizePostQueryURLAndCompareResponse(t *testing.T, server *httptest.Server, queryURL *url.URL, redirectURL *url.URL) {
	res, err := server.Client().Post(queryURL.String(), "multipart/form-data", bytes.NewBuffer([]byte("")))
	require.NoError(t, err)
	require.Equal(t, 200, res.StatusCode)
	require.Equal(t, string(templateMustExec(t, templateAuthTokenServerRedirect, map[string]template.URL{
		"URL": template.URL(redirectURL.String()),
	})), string(mustReadAllBytes(t, res.Body)))
}

func templateMustExec(t *testing.T, tpl *template.Template, args interface{}) []byte {
	b := bytes.NewBuffer(nil)
	err := tpl.Execute(b, args)
	require.NoError(t, err)

	return b.Bytes()
}

func setURLParam(u *url.URL, k, v string) {
	q := u.Query()
	q.Set(k, v)
	u.RawQuery = q.Encode()
}
