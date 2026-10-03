package orm_test

import (
	"context"
	"strings"
	"testing"

	"sumeru/core/orm"
)

func TestFormatAuditLogBodyWriteDiff(t *testing.T) {
	ctx := context.Background()
	before := `{"name":"Alice","active":true}`
	after := `{"name":"Bob","active":true}`
	body := orm.FormatAuditLogBodyForTest(ctx, "core.user", "write", before, after, "")
	if !strings.Contains(body, "→") && !strings.Contains(body, "Bob") {
		t.Fatalf("expected field diff, got %q", body)
	}
}

func TestFormatAuditLogBodyUnlink(t *testing.T) {
	body := orm.FormatAuditLogBodyForTest(context.Background(), "core.user", "unlink", "", "", "")
	if body != "Deleted record" {
		t.Fatalf("got %q", body)
	}
}

func TestFormatAuditLogBodyCreate(t *testing.T) {
	after := `{"name":"Acme"}`
	body := orm.FormatAuditLogBodyForTest(context.Background(), "core.company", "create", "", after, "")
	if body == "" || body == "Change recorded" {
		t.Fatalf("expected create body, got %q", body)
	}
}

func TestFormatAuditLogBodyRedactsSensitive(t *testing.T) {
	ctx := context.Background()
	before := `{"password":"old"}`
	after := `{"password":"new"}`
	body := orm.FormatAuditLogBodyForTest(ctx, "core.user", "write", before, after, "")
	if strings.Contains(strings.ToLower(body), "password") {
		t.Fatalf("expected redacted field, got %q", body)
	}
}
