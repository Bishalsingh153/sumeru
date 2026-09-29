package module_test

import (
	"strings"
	"testing"

	"sumeru/core/module"
)

func TestRegistryHasField_unknownModel(t *testing.T) {
	if module.RegistryHasFieldForTest("not.a.model", "login") {
		t.Fatal("expected unknown model")
	}
}

func TestFieldAccessDefaultName(t *testing.T) {
	name := module.FieldAccessDefaultNameForTest("core.user", "password", "base.group_system")
	if !strings.Contains(name, "core.user") || !strings.Contains(name, "password") {
		t.Fatalf("unexpected default name: %q", name)
	}
	global := module.FieldAccessDefaultNameForTest("core.user", "password", "")
	if !strings.Contains(global, "global") {
		t.Fatalf("expected global segment: %q", global)
	}
}
