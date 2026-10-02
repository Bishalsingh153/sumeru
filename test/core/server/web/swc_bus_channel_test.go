package web_test

import (
	"context"
	"testing"

	"sumeru/core/server/web"
)

func TestAuthorizeSwcBusChannel_userScope(t *testing.T) {
	ctx := context.Background()
	if web.AuthorizeSwcBusChannel(ctx, 5, "user/5/notifications") != true {
		t.Fatal("expected own user channel allowed")
	}
	if web.AuthorizeSwcBusChannel(ctx, 5, "user/6/notifications") {
		t.Fatal("expected other user channel denied")
	}
}

func TestRecordBusChannel_format(t *testing.T) {
	ch := web.RecordBusChannel("crm.lead", 9)
	if ch != "record/crm.lead/9" {
		t.Fatalf("channel: %q", ch)
	}
}
