package services

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/oauth2"
)

// ghRoundTripper intercepts calls to api.github.com so the GitHub-dependent OAuth
// paths can be exercised without reaching the real API. Any other host (the mock
// token server) is delegated to the default transport.
type ghRoundTripper struct {
	userStatus  int
	userBody    string
	emailStatus int
	emailBody   string
}

func (rt *ghRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "api.github.com" {
		switch req.URL.Path {
		case "/user":
			return jsonResponse(rt.userStatus, rt.userBody), nil
		case "/user/emails":
			return jsonResponse(rt.emailStatus, rt.emailBody), nil
		}
	}
	return http.DefaultTransport.RoundTrip(req)
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}
}

func newMockOAuthService(tokenServerURL string, rt http.RoundTripper) (*OAuthService, context.Context) {
	svc := &OAuthService{
		userRepo: newMockRepo(),
		oauthConfig: &oauth2.Config{
			ClientID:     "client-id",
			ClientSecret: "client-secret",
			Endpoint: oauth2.Endpoint{
				AuthURL:  tokenServerURL + "/authorize",
				TokenURL: tokenServerURL + "/token",
			},
		},
	}
	ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{Transport: rt})
	return svc, ctx
}

func tokenServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, `{"access_token":"mock-access-token","token_type":"bearer"}`)
	}))
}

func TestExchangeCodeForToken_Success(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	svc, ctx := newMockOAuthService(ts.URL, &ghRoundTripper{})
	tok, err := svc.ExchangeCodeForToken(ctx, "valid-code")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.AccessToken != "mock-access-token" {
		t.Fatalf("expected mock token, got %q", tok.AccessToken)
	}
}

func TestExchangeCodeForToken_EmptyCode(t *testing.T) {
	svc, ctx := newMockOAuthService("http://unused", &ghRoundTripper{})
	if _, err := svc.ExchangeCodeForToken(ctx, ""); err == nil {
		t.Fatal("expected error for empty code")
	}
}

func TestFetchGitHubUser_WithEmail(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	rt := &ghRoundTripper{
		userStatus: 200,
		userBody:   `{"id":99,"login":"octocat","name":"The Octocat","email":"octo@github.com","avatar_url":"http://a"}`,
	}
	svc, ctx := newMockOAuthService(ts.URL, rt)

	user, err := svc.FetchGitHubUser(ctx, &oauth2.Token{AccessToken: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "octo@github.com" || user.Login != "octocat" {
		t.Fatalf("unexpected user: %+v", user)
	}
}

func TestFetchGitHubUser_FallsBackToPrimaryEmail(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	rt := &ghRoundTripper{
		userStatus:  200,
		userBody:    `{"id":99,"login":"octocat","email":""}`,
		emailStatus: 200,
		emailBody:   `[{"email":"secondary@x.com","primary":false,"verified":true},{"email":"primary@x.com","primary":true,"verified":true}]`,
	}
	svc, ctx := newMockOAuthService(ts.URL, rt)

	user, err := svc.FetchGitHubUser(ctx, &oauth2.Token{AccessToken: "x"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "primary@x.com" {
		t.Fatalf("expected primary verified email, got %q", user.Email)
	}
}

func TestFetchGitHubUser_UserEndpointError(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	rt := &ghRoundTripper{userStatus: 401, userBody: `{"message":"Bad credentials"}`}
	svc, ctx := newMockOAuthService(ts.URL, rt)

	if _, err := svc.FetchGitHubUser(ctx, &oauth2.Token{AccessToken: "x"}); err == nil {
		t.Fatal("expected error when github /user returns non-200")
	}
}

func TestFetchGitHubUser_NoVerifiedPrimaryEmail(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	rt := &ghRoundTripper{
		userStatus:  200,
		userBody:    `{"id":99,"login":"octocat","email":""}`,
		emailStatus: 200,
		emailBody:   `[{"email":"unverified@x.com","primary":true,"verified":false}]`,
	}
	svc, ctx := newMockOAuthService(ts.URL, rt)

	if _, err := svc.FetchGitHubUser(ctx, &oauth2.Token{AccessToken: "x"}); err == nil {
		t.Fatal("expected error when no verified primary email exists")
	}
}

func TestFetchGitHubUser_EmailEndpointError(t *testing.T) {
	ts := tokenServer(t)
	defer ts.Close()

	rt := &ghRoundTripper{
		userStatus:  200,
		userBody:    `{"id":99,"login":"octocat","email":""}`,
		emailStatus: 500,
		emailBody:   `{}`,
	}
	svc, ctx := newMockOAuthService(ts.URL, rt)

	if _, err := svc.FetchGitHubUser(ctx, &oauth2.Token{AccessToken: "x"}); err == nil {
		t.Fatal("expected error when github /user/emails returns non-200")
	}
}
