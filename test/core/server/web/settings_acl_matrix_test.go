package web_test

import (
	"testing"

	"sumeru/core/server/web"
)

func TestFieldACLMatrixRowName(t *testing.T) {
	global := web.FieldACLMatrixRowNameForTest("core.user", "password", 0)
	if global != "matrix.core.user.password.global" {
		t.Fatalf("global name = %q", global)
	}
	grouped := web.FieldACLMatrixRowNameForTest("core.user", "password", 3)
	if grouped != "matrix.core.user.password.g3" {
		t.Fatalf("group name = %q", grouped)
	}
}

func TestModelACLMatrixRowName(t *testing.T) {
	if got := web.ModelACLMatrixRowNameForTest("core.user", 0); got != "matrix.core.user.global" {
		t.Fatalf("global = %q", got)
	}
	if got := web.ModelACLMatrixRowNameForTest("core.user", 5); got != "matrix.core.user.g5" {
		t.Fatalf("group = %q", got)
	}
}

func TestMatrixGroupsForDisplay_truncates(t *testing.T) {
	displayLen, truncated := web.MatrixGroupDisplayLenForTest(50, false)
	if !truncated {
		t.Fatal("expected truncation")
	}
	if displayLen != 40 {
		t.Fatalf("display len = %d; want 40 (includes Global column)", displayLen)
	}
}

func TestFilterFieldNamesForTest(t *testing.T) {
	fields := []string{"email", "login", "name", "password"}
	got := web.FilterFieldNamesForTest(fields, "pass")
	if len(got) != 1 || got[0] != "password" {
		t.Fatalf("filter = %v", got)
	}
}
