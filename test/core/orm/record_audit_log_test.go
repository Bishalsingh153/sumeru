package orm_test

import (
	"context"
	"testing"

	"sumeru/core/orm"
)

func TestSearchRecordAuditLog_invalidRecord(t *testing.T) {
	_, err := orm.SearchRecordAuditLogForTest(context.Background(), "", 0, 10, 0)
	if err == nil {
		t.Fatal("expected error for empty model")
	}
}

func TestRecordAuditLogAvailable_withoutAuditModel(t *testing.T) {
	if _, ok := orm.Registry["sys.audit"]; ok {
		t.Skip("sys.audit registered in this test run")
	}
	if orm.RecordAuditLogAvailableForTest(context.Background()) {
		t.Fatal("expected false without sys.audit model")
	}
}
