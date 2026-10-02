package web_test

import (
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"sumeru/core/engine/render"
	"sumeru/core/server/config"
	"sumeru/core/server/web"
)

func TestBuildLoginBrand_logoPriority(t *testing.T) {
	render.SetShellBranding(render.ShellBranding{})
	t.Cleanup(func() { render.SetShellBranding(render.ShellBranding{}) })

	defaultBrand := web.BuildLoginBrandForTest(web.LoginCompanyRowForTest{})
	if defaultBrand.LogoURL != web.TestDefaultSumeruLogoURL {
		t.Fatalf("default logo: got %q want %q", defaultBrand.LogoURL, web.TestDefaultSumeruLogoURL)
	}
	if defaultBrand.ShowPoweredBy {
		t.Fatal("default brand should not show powered-by")
	}
	if !defaultBrand.ShowTrustChips {
		t.Fatal("default brand should show trust chips")
	}

	withCompanyLogo := web.BuildLoginBrandForTest(web.LoginCompanyRowForTest{
		Name:      "Acme",
		LoginLogo: "data:image/png;base64,abcd",
	})
	if withCompanyLogo.LogoURL != web.TestLoginBrandLogoRoute {
		t.Fatalf("company logo route: got %q", withCompanyLogo.LogoURL)
	}
	if !withCompanyLogo.ShowPoweredBy {
		t.Fatal("tenant brand should show powered-by")
	}

	render.SetShellBranding(render.ShellBranding{LogoURL: "/static/app-logo"})
	cfgBrand := web.BuildLoginBrandForTest(web.LoginCompanyRowForTest{})
	if cfgBrand.LogoURL != "/static/app-logo" {
		t.Fatalf("config logo: got %q", cfgBrand.LogoURL)
	}
}

func TestParseLoginLogoDataURL(t *testing.T) {
	png := []byte{0x89, 0x50, 0x4e, 0x47}
	raw := "data:image/png;base64," + base64.StdEncoding.EncodeToString(png)
	mime, data, ok := web.ParseLoginLogoDataURLForTest(raw)
	if !ok || mime != "image/png" || len(data) != len(png) {
		t.Fatalf("parse ok=%v mime=%q len=%d", ok, mime, len(data))
	}
	if _, _, bad := web.ParseLoginLogoDataURLForTest("data:text/plain;base64,YQ=="); bad {
		t.Fatal("expected reject non-image mime")
	}
}

func TestLoginBrandLogoGet_notFoundWithoutLogo(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, web.TestLoginBrandLogoRoute, nil)
	rec := httptest.NewRecorder()
	web.LoginBrandLogoGetForTest(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d", rec.Code)
	}
}

func TestLoginGet_rendersBrandLockup(t *testing.T) {
	root := sumeruModuleRoot(t)
	prevTemplates := config.AppConfig.TemplatesPath
	config.AppConfig.TemplatesPath = filepath.Join(root, "core", "engine", "templates")
	t.Cleanup(func() { config.AppConfig.TemplatesPath = prevTemplates })
	web.ResetAuthTemplateCacheForTest()

	req := httptest.NewRequest(http.MethodGet, web.TestLoginRoute, nil)
	rec := httptest.NewRecorder()
	web.LoginGetForTest(rec, req)
	body := rec.Body.String()
	for _, needle := range []string{
		"sum-login-brand-lockup",
		web.TestDefaultSumeruLogoURL,
		"sum-login-trust-chip",
	} {
		if !strings.Contains(body, needle) {
			t.Fatalf("missing %q in login HTML", needle)
		}
	}
}
