package orm_test

import (
	"context"
	"net/url"
	"strings"
	"testing"

	"sumeru/core/orm"
)

func TestApplyUserSecurityPostErr_passwordMismatch(t *testing.T) {
	ctx := orm.ContextWithUID(context.Background(), 1)
	form := url.Values{}
	form.Set("password_plain", "abc")
	form.Set("password_plain_confirm", "xyz")
	err := orm.ApplyUserSecurityPostErr(ctx, 1, 1, form)
	if err == nil || !strings.Contains(err.Error(), "match") {
		t.Fatalf("expected mismatch error, got %v", err)
	}
}
