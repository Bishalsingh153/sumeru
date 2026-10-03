package module_test

import (
	"context"
	"testing"

	"sumeru/core/module"
)

func TestExpandInstallModuleNamesExplicitOrder(t *testing.T) {
	module.DiscoveredAddons = map[string]*module.Addon{
		"base": {Manifest: module.Manifest{Name: "base", Depends: []string{}}},
		"im":  {Manifest: module.Manifest{Name: "im", Depends: []string{"base"}}},
		"crm": {Manifest: module.Manifest{Name: "crm", Depends: []string{"base", "im"}}},
	}
	ctx := context.Background()
	names, err := module.ExpandInstallModuleNamesForTest(ctx, []string{"crm", "im"})
	if err != nil {
		t.Fatal(err)
	}
	if len(names) != 2 {
		t.Fatalf("expected 2 modules, got %v", names)
	}
	if names[0] != "im" || names[1] != "crm" {
		t.Fatalf("unexpected order: %v", names)
	}
}
