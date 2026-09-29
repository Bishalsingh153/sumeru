package security_test

import (
	"testing"

	"sumeru/core/security"
)

func TestFieldKernelPolicyTags_password(t *testing.T) {
	tags := security.FieldKernelPolicyTags("core.user", "password")
	if len(tags) < 2 {
		t.Fatalf("expected kernel tags for password; got %v", tags)
	}
}

func TestFieldKernelPolicyTags_unknownField(t *testing.T) {
	if tags := security.FieldKernelPolicyTags("core.user", "login"); len(tags) != 0 {
		t.Fatalf("login should have no kernel tags; got %v", tags)
	}
}
