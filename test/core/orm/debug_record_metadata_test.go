package orm_test

import (
	"context"
	"strings"
	"testing"

	"sumeru/core/orm"
)

func TestBuildDebugRecordMetadata_requiresLogin(t *testing.T) {
	_, err := orm.BuildDebugRecordMetadata(context.Background(), 0, "core.partner", 1)
	if err == nil || !strings.Contains(err.Error(), "login") {
		t.Fatalf("expected login error, got %v", err)
	}
}

func TestBuildDebugRecordMetadata_requiresRecordID(t *testing.T) {
	_, err := orm.BuildDebugRecordMetadata(context.Background(), 2, "core.partner", 0)
	if err == nil || !strings.Contains(err.Error(), "record_id") {
		t.Fatalf("expected record_id error, got %v", err)
	}
}
