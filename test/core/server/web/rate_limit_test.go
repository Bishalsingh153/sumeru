package web_test

import (
	"testing"

	"sumeru/core/server/web"
)

func TestRateLimitedPath_authRoutes(t *testing.T) {
	for _, path := range []string{
		web.TestLoginRoute,
		"/web/login/totp",
		"/web/auth/oauth/start",
		"/web/auth/oauth/callback",
		web.TestAPIRPCRoute,
	} {
		if !web.RateLimitedPathForTest(path) {
			t.Fatalf("expected rate limit for %q", path)
		}
	}
	if web.RateLimitedPathForTest("/web/home") {
		t.Fatal("home should not be rate limited")
	}
}
