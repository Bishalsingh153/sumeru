package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sumeru/core/server/web"
)

func TestOAuthCallback_invalidStateRedirectsToLogin(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/web/auth/oauth/callback?state=bad&code=abc", nil)
	rec := httptest.NewRecorder()
	web.OAuthCallbackForTest(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("status=%d", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if !strings.Contains(loc, web.TestOAuthDeniedMsg) {
		t.Fatalf("location=%q", loc)
	}
}

func TestOAuthStart_missingProvider(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/web/auth/oauth/start", nil)
	rec := httptest.NewRecorder()
	web.OAuthStartForTest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestAuthLocalEnabled_defaultsTrue(t *testing.T) {
	if !web.AuthLocalEnabledForTest(context.Background()) {
		t.Fatal("expected local login enabled by default")
	}
}

func TestListEnabledAuthProviders_emptyWithoutDB(t *testing.T) {
	if got := web.ListEnabledAuthProvidersForTest(context.Background()); len(got) != 0 {
		t.Fatalf("expected no providers, got %v", got)
	}
	_ = web.LoginPageCompanyNameForTest(context.Background())
}
