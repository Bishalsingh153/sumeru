package orm_test

import (
	"context"
	"testing"

	"sumeru/core/orm"
)

func TestResolveContentAttachmentInvalid(t *testing.T) {
	_, err := orm.ResolveContentAttachmentForTest(context.Background(), 0, "", "", 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestResolveContentAttachmentByIDWithoutDB(t *testing.T) {
	_, err := orm.ResolveContentAttachmentForTest(context.Background(), 99, "", "", 0)
	if err == nil {
		t.Fatal("expected error without db")
	}
}

func TestCanReadAttachmentContentNoAccess(t *testing.T) {
	err := orm.CanReadAttachmentContentForTest(context.Background(), 0, map[string]interface{}{})
	if err == nil {
		t.Fatal("expected access error")
	}
}
