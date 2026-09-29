package render_test

import (
	"testing"
	"sumeru/core/engine/parser"
	"sumeru/core/engine/render"
)


func TestBuildSidebarMenus_localizationSection(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "160", Name: "Localization", ParentID: "100", Sequence: 60, AccessGroups: "base.group_system"},
		{ID: "161", Name: "Countries", ParentID: "160", Sequence: 10, Action: "/web?action=10&menu_id=161", AccessGroups: "base.group_system"},
		{ID: "162", Name: "States", ParentID: "160", Sequence: 20, Action: "/web?action=11&menu_id=162", AccessGroups: "base.group_system"},
		{ID: "163", Name: "Cities", ParentID: "160", Sequence: 30, Action: "/web?action=12&menu_id=163", AccessGroups: "base.group_system"},
	}
	menuAllowed := func(mi parser.MenuItem) bool {
		return mi.AccessGroups == "base.group_system"
	}
	sections := render.BuildSidebarMenus(allMenus, "100", menuAllowed)
	var localization *render.SidebarMenu
	for i := range sections {
		if sections[i].Name == "Localization" {
			localization = &sections[i]
			break
		}
	}
	if localization == nil {
		t.Fatal("Localization section not found in Settings sidebar")
	}
	if len(localization.SubMenus) != 3 {
		t.Fatalf("Localization submenus = %d; want 3", len(localization.SubMenus))
	}
	names := map[string]bool{}
	for _, sm := range localization.SubMenus {
		names[sm.Name] = true
		if sm.Action == "" {
			t.Fatalf("submenu %q has empty action URL", sm.Name)
		}
	}
	for _, want := range []string{"Countries", "States", "Cities"} {
		if !names[want] {
			t.Fatalf("missing submenu %q", want)
		}
	}
}

func TestBuildSidebarMenus_companiesConfigurationFlatten(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "110", Name: "Companies", ParentID: "100", Sequence: 10},
		{ID: "111", Name: "All Companies", ParentID: "110", Sequence: 10, Action: "/web?action=1&menu_id=111"},
		{ID: "112", Name: "Configuration", ParentID: "110", Sequence: 90},
		{ID: "113", Name: "Access Rights", ParentID: "112", Sequence: 16, Action: "/web?action=2&menu_id=113"},
		{ID: "114", Name: "Field access matrix", ParentID: "112", Sequence: 19, Action: "/web?action=3&menu_id=114"},
	}
	menuAllowed := func(parser.MenuItem) bool { return true }
	sections := render.BuildSidebarMenus(allMenus, "100", menuAllowed)
	var companies *render.SidebarMenu
	for i := range sections {
		if sections[i].Name == "Companies" {
			companies = &sections[i]
			break
		}
	}
	if companies == nil {
		t.Fatal("Companies section not found")
	}
	names := map[string]bool{}
	for _, sm := range companies.SubMenus {
		names[sm.Name] = true
	}
	if names["Configuration"] {
		t.Fatal("Configuration container should not appear as a link")
	}
	for _, want := range []string{"All Companies", "Access Rights", "Field access matrix"} {
		if !names[want] {
			t.Fatalf("missing flattened link %q; got %v", want, names)
		}
	}
}

func TestApplySettingsNavHrefOverridesForTest_staleWindowAction(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "101", Name: "General Settings", ParentID: "100", Sequence: 1, Action: "/web?action=99&menu_id=101"},
	}
	render.ApplySettingsNavHrefOverridesForTest(allMenus, "101")
	if allMenus[1].Action != render.SettingsHubPath {
		t.Fatalf("Action = %q; want %q", allMenus[1].Action, render.SettingsHubPath)
	}
	sections := render.BuildSidebarMenus(allMenus, "100", func(parser.MenuItem) bool { return true })
	if len(sections) != 1 || sections[0].SubMenus[0].Action != render.SettingsHubPath {
		t.Fatalf("sidebar after override: %+v", sections)
	}
}

func TestBuildSidebarMenus_generalSettingsSection(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "101", Name: "General Settings", ParentID: "100", Sequence: 1, Action: "/web/settings"},
	}
	sections := render.BuildSidebarMenus(allMenus, "100", func(parser.MenuItem) bool { return true })
	if len(sections) != 1 {
		t.Fatalf("sections = %d; want 1", len(sections))
	}
	if sections[0].Name != "General Settings" || len(sections[0].SubMenus) != 1 {
		t.Fatalf("unexpected section: %+v", sections[0])
	}
	if sections[0].SubMenus[0].Name != "General Settings" || sections[0].SubMenus[0].Action != "/web/settings" {
		t.Fatalf("general settings link = %+v", sections[0].SubMenus[0])
	}
}

func TestBuildSidebarMenus_excludesPersonalAccountSecurity(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "100", Name: "Settings", Sequence: 5},
		{ID: "101", Name: "General Settings", ParentID: "100", Sequence: 1, Action: "/web/settings"},
		{ID: "120", Name: "Personal", ParentID: "100", Sequence: 2},
		{ID: "121", Name: "Account security", ParentID: "120", Sequence: 10, Action: "/web/settings/account"},
		{ID: "110", Name: "Companies", ParentID: "100", Sequence: 10},
		{ID: "111", Name: "All Companies", ParentID: "110", Sequence: 10, Action: "/web?action=1&menu_id=111"},
	}
	sections := render.BuildSidebarMenus(allMenus, "100", func(parser.MenuItem) bool { return true })
	for _, sec := range sections {
		if sec.Name == "Personal" {
			t.Fatal("Personal section must not appear in settings sidebar")
		}
		for _, sm := range sec.SubMenus {
			if sm.Name == "Account security" || sm.Action == "/web/settings/account" {
				t.Fatalf("account security must not appear in sidebar: section %q link %q", sec.Name, sm.Name)
			}
		}
	}
}

func TestBuildSidebarMenus_skipsEmptySections(t *testing.T) {
	allMenus := []parser.MenuItem{
		{ID: "200", Name: "Contacts", Sequence: 20},
		{ID: "210", Name: "Contacts organization", ParentID: "200", Sequence: 100},
		{ID: "211", Name: "All Contacts", ParentID: "210", Sequence: 10, Action: "/web?action=1&menu_id=211"},
		{ID: "220", Name: "Empty section", ParentID: "200", Sequence: 200},
	}
	menuAllowed := func(parser.MenuItem) bool { return true }
	sections := render.BuildSidebarMenus(allMenus, "200", menuAllowed)
	if len(sections) != 1 {
		t.Fatalf("sections = %d; want 1 (empty section skipped)", len(sections))
	}
	if sections[0].Name != "Contacts organization" {
		t.Fatalf("section name = %q; want Contacts organization", sections[0].Name)
	}
	if len(sections[0].SubMenus) != 1 || sections[0].SubMenus[0].Name != "All Contacts" {
		t.Fatalf("nested contact link: %+v", sections[0].SubMenus)
	}
}

func TestSidebarHasMenus(t *testing.T) {
	if render.SidebarHasMenus(nil) {
		t.Fatal("nil menus should be false")
	}
	if render.SidebarHasMenus([]render.SidebarMenu{{Name: "Empty", SubMenus: nil}}) {
		t.Fatal("section without links should be false")
	}
	if !render.SidebarHasMenus([]render.SidebarMenu{{Name: "Links", SubMenus: []parser.MenuItem{{Name: "One"}}}}) {
		t.Fatal("section with links should be true")
	}
}

func TestResolveActiveModuleID_noDefaultWhenMenuMissing(t *testing.T) {
	got := render.ResolveActiveModuleID(nil, "")
	if got != "" {
		t.Fatalf("active module = %q; want empty when menu_id missing", got)
	}
}
