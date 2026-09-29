package web_test

import (
	"testing"

	"sumeru/core/server/web"
)

func TestBuildSettingsHubPageData(t *testing.T) {
	page := web.BuildSettingsHubPageData(t.Context(), "5")
	if page.Title != web.TestSettingsHubPageTitle || !page.SettingsNavActive || page.ActiveMenuID != "5" {
		t.Fatalf("unexpected page data: %+v", page)
	}
	if page.SuppressSidebar {
		t.Fatal("settings hub should show the settings sidebar")
	}
	if page.SuppressActivityDock {
		t.Fatal("settings hub should show the activity panel")
	}
	if len(page.ViewStylesheetURLs) != 1 || page.ViewStylesheetURLs[0] != web.TestSettingsHubStylesheetURL {
		t.Fatalf("unexpected stylesheets: %v", page.ViewStylesheetURLs)
	}
	if len(page.ExtraScriptURLs) != 1 || page.ExtraScriptURLs[0] != web.TestSettingsHubScriptURL {
		t.Fatalf("unexpected scripts: %v", page.ExtraScriptURLs)
	}
}
