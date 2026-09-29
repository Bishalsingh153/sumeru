package render_test

import (
	"testing"

	"sumeru/core/engine/parser"
	"sumeru/core/engine/render"
)

func TestBuildSettingsHubCategories_securityAndDedupe(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "105", Name: "General Settings", ParentID: "100", Sequence: 1, Action: "/web?action=99&menu_id=105"},
		{ID: "115", Name: "Security", ParentID: "100", Sequence: 15},
		{ID: "116", Name: "Access Rights", ParentID: "115", Sequence: 10, Action: "/web?action=1&menu_id=116"},
		{ID: "117", Name: "Field access matrix", ParentID: "115", Sequence: 40, Action: "/web?action=2&menu_id=117"},
		{ID: "110", Name: "Companies", ParentID: "100", Sequence: 10},
		{ID: "111", Name: "All Companies", ParentID: "110", Sequence: 10, Action: "/web?action=9&menu_id=111"},
		{ID: "112", Name: "Duplicate Companies", ParentID: "110", Sequence: 20, Action: "/web?action=9&menu_id=112"},
	}
	menuAllowed := func(parser.MenuItem) bool { return true }
	categories := render.BuildSettingsHubCategoriesForTest(allMenus, "100", menuAllowed, "105")
	if len(categories) < 2 {
		t.Fatalf("categories = %d; want at least 2", len(categories))
	}
	for _, c := range categories {
		if c.Title == "General Settings" || c.Title == "Personal" {
			t.Fatalf("hub nav menu must not appear as a category card: %q", c.Title)
		}
	}
	var security *render.SettingsHubCategory
	for i := range categories {
		if categories[i].Title == "Security" {
			security = &categories[i]
			break
		}
	}
	if security == nil {
		t.Fatal("Security category missing")
	}
	linkNames := map[string]bool{}
	for _, g := range security.Groups {
		for _, l := range g.Links {
			linkNames[l.Name] = true
		}
	}
	if !linkNames["Field access matrix"] || !linkNames["Access Rights"] {
		t.Fatalf("security links: %v", linkNames)
	}
	var companies *render.SettingsHubCategory
	for i := range categories {
		if categories[i].Title == "Companies" {
			companies = &categories[i]
			break
		}
	}
	if companies == nil {
		t.Fatal("Companies category missing")
	}
	companyLinks := 0
	for _, g := range companies.Groups {
		companyLinks += len(g.Links)
	}
	if companyLinks != 1 {
		t.Fatalf("company links = %d; want 1 (deduped action)", companyLinks)
	}
}

func TestBuildSettingsHubCategories_excludesPersonalAccountSecurity(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "105", Name: "General Settings", ParentID: "100", Sequence: 1, Action: "/web?action=99&menu_id=105"},
		{ID: "120", Name: "Personal", ParentID: "100", Sequence: 2},
		{ID: "121", Name: "Account security", ParentID: "120", Sequence: 10, Action: "/web/settings/account"},
		{ID: "115", Name: "Security", ParentID: "100", Sequence: 15},
		{ID: "116", Name: "Access Rights", ParentID: "115", Sequence: 10, Action: "/web?action=1&menu_id=116"},
	}
	menuAllowed := func(parser.MenuItem) bool { return true }
	categories := render.BuildSettingsHubCategoriesForTest(allMenus, "100", menuAllowed, "105")
	for _, c := range categories {
		if c.Title == "Personal" {
			t.Fatal("Personal must not appear on the settings hub")
		}
		for _, g := range c.Groups {
			for _, l := range g.Links {
				if l.Name == "Account security" || l.Href == "/web/settings/account" {
					t.Fatalf("account security link must not appear on hub: category %q link %q", c.Title, l.Name)
				}
			}
		}
	}
}
