package api

import (
	"context"
	"mayfly-go/pkg/req"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestOperatorAllowlist(t *testing.T) {
	if operatorAllowed(0, "0") || operatorAllowed(2, "1") || operatorAllowed(1, "") {
		t.Fatal("unauthorized operator accepted")
	}
	if !operatorAllowed(1, " 1, 2 ") {
		t.Fatal("authorized operator rejected")
	}
}
func TestBridgeRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/admin/accounts" || r.Header.Get("X-RAYS-Admin-Actor") != "mayfly:1" || r.Header.Get("X-RAYS-Admin-Token") != strings.Repeat("s", 32) {
			t.Error("invalid upstream request")
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success":true,"data":[{"id":"test"}]}`))
	}))
	defer server.Close()
	t.Setenv("RAYS_ADMIN_API_URL", server.URL)
	t.Setenv("RAYS_ADMIN_BRIDGE_TOKEN", strings.Repeat("s", 32))
	result, err := bridgeRequest(context.Background(), "GET", "/accounts", "", strings.NewReader(""), 1)
	if err != nil || !strings.Contains(string(result), "test") {
		t.Fatalf("unexpected response: %v", err)
	}
	t.Setenv("RAYS_ADMIN_BRIDGE_TOKEN", "")
	if _, err = bridgeRequest(context.Background(), "GET", "/accounts", "", strings.NewReader(""), 1); err == nil {
		t.Fatal("missing credential accepted")
	}
}
func TestBridgeDoesNotFollowRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { http.Redirect(w, r, "http://invalid.example", 302) }))
	defer server.Close()
	t.Setenv("RAYS_ADMIN_API_URL", server.URL)
	t.Setenv("RAYS_ADMIN_BRIDGE_TOKEN", strings.Repeat("s", 32))
	if _, err := bridgeRequest(context.Background(), "GET", "/overview", "", strings.NewReader(""), 1); err == nil {
		t.Fatal("redirect accepted")
	}
}

func TestLogoutRevokesAccessAndRefreshTokens(t *testing.T) {
	access, refresh, err := req.CreateToken(901, "qa")
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseAccessToken(access); err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseToken(refresh); err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseAccessToken(refresh); err == nil {
		t.Fatal("refresh token accepted as access")
	}
	if err = req.RevokeTokens(901); err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseAccessToken(access); err == nil {
		t.Fatal("logged-out access token accepted")
	}
	if _, _, err = req.ParseToken(refresh); err == nil {
		t.Fatal("logged-out refresh token accepted")
	}
	next, _, err := req.CreateToken(901, "qa")
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseAccessToken(next); err != nil {
		t.Fatal(err)
	}
	if _, _, err = req.ParseAccessToken(access); err == nil {
		t.Fatal("new login resurrected old token")
	}
}
