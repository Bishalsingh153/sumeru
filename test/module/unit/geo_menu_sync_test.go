package module_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sumeru/core/engine/parser"
	"sumeru/core/module"
)

func addonDir(t *testing.T, name string) string {
	t.Helper()
	dir := filepath.Join("..", "..", "..", "addons", name)
	if _, err := os.Stat(dir); err != nil {
		t.Skip(name, " addon not found at ", dir)
	}
	return dir
}

func menuByID(menus []parser.MenuItem, id string) (parser.MenuItem, bool) {
	for _, m := range menus {
		if m.ID == id {
			return m, true
		}
	}
	return parser.MenuItem{}, false
}

func collectManifestMenus(t *testing.T, addonDir string) []parser.MenuItem {
	t.Helper()
	manifestPath := filepath.Join(addonDir, "manifest.json")
	raw, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatalf("read manifest: %v", err)
	}
	var manifest struct {
		Data []string `json:"data"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		t.Fatalf("parse manifest: %v", err)
	}
	var out []parser.MenuItem
	for _, rel := range manifest.Data {
		if !strings.HasSuffix(strings.ToLower(rel), ".xml") {
			continue
		}
		items, err := module.CollectMenuItemsFromManifestFile(filepath.Join(addonDir, rel))
		if err != nil {
			t.Fatalf("collect menus from %s: %v", rel, err)
		}
		out = append(out, items...)
	}
	return out
}

func TestGeoDeferredMenusIncludeGeoSection(t *testing.T) {
	menus := collectManifestMenus(t, addonDir(t, "geo"))
	geoSection, ok := menuByID(menus, "menu_geo_section")
	if !ok {
		t.Fatal("menu_geo_section missing from geo manifest menus")
	}
	if geoSection.ParentID != "base.menu_settings_root" {
		t.Fatalf("menu_geo_section parent = %q; want base.menu_settings_root", geoSection.ParentID)
	}
	for _, childID := range []string{"menu_core_country", "menu_core_country_state", "menu_core_city"} {
		child, ok := menuByID(menus, childID)
		if !ok {
			t.Fatalf("%s missing", childID)
		}
		if child.ParentID != "menu_geo_section" {
			t.Fatalf("%s parent = %q; want menu_geo_section", childID, child.ParentID)
		}
	}
}

func TestI18nDeferredMenusIncludeI18nSection(t *testing.T) {
	menus := collectManifestMenus(t, addonDir(t, "i18n"))
	section, ok := menuByID(menus, "menu_i18n_section")
	if !ok {
		t.Fatal("menu_i18n_section missing")
	}
	if section.Name != "Internationalization" {
		t.Fatalf("section name = %q", section.Name)
	}
	for _, childID := range []string{"menu_core_lang", "menu_sys_translation"} {
		if _, ok := menuByID(menus, childID); !ok {
			t.Fatalf("%s missing", childID)
		}
	}
}

func TestGeoManifestMenusLoadLast(t *testing.T) {
	dir := addonDir(t, "geo")
	raw, _ := os.ReadFile(filepath.Join(dir, "manifest.json"))
	var manifest struct {
		Data []string `json:"data"`
	}
	_ = json.Unmarshal(raw, &manifest)
	last := manifest.Data[len(manifest.Data)-1]
	if last != "views/menus.xml" {
		t.Fatalf("geo menus.xml must be last; got %q", last)
	}
}

func TestI18nManifestMenusLoadLast(t *testing.T) {
	dir := addonDir(t, "i18n")
	raw, _ := os.ReadFile(filepath.Join(dir, "manifest.json"))
	var manifest struct {
		Data []string `json:"data"`
	}
	_ = json.Unmarshal(raw, &manifest)
	last := manifest.Data[len(manifest.Data)-1]
	if last != "views/menus.xml" {
		t.Fatalf("i18n menus.xml must be last; got %q", last)
	}
}

func TestBaseManifestMenusLoadLast(t *testing.T) {
	dir := addonDir(t, "base")
	raw, _ := os.ReadFile(filepath.Join(dir, "manifest.json"))
	var manifest struct {
		Data []string `json:"data"`
	}
	_ = json.Unmarshal(raw, &manifest)
	last := manifest.Data[len(manifest.Data)-1]
	if last != "views/menus.xml" {
		t.Fatalf("base menus.xml must be last; got %q", last)
	}
}
