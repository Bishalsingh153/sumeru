package audit_test

import (
	"testing"

	"sumeru/addons/audit"
)

func TestRetentionDaysFromRecord(t *testing.T) {
	d, err := audit.RetentionDaysFromRecordForTest(map[string]interface{}{"retention_preset": "7_days"})
	if err != nil || d != 7 {
		t.Fatalf("7_days: d=%d err=%v", d, err)
	}
	d, err = audit.RetentionDaysFromRecordForTest(map[string]interface{}{"retention_preset": "custom", "retention_days": 14})
	if err != nil || d != 14 {
		t.Fatalf("custom: d=%d err=%v", d, err)
	}
	_, err = audit.RetentionDaysFromRecordForTest(map[string]interface{}{"retention_preset": "custom"})
	if err == nil {
		t.Fatal("expected custom days error")
	}
	d, err = audit.RetentionDaysFromRecordForTest(map[string]interface{}{"retention_preset": "30_days"})
	if err != nil || d != 30 {
		t.Fatalf("30_days: d=%d err=%v", d, err)
	}
	_, err = audit.RetentionDaysFromRecordForTest(map[string]interface{}{"retention_preset": "never"})
	if err == nil {
		t.Fatal("expected unknown preset error")
	}
	_, err = audit.RetentionDaysFromRecordForTest(map[string]interface{}{
		"retention_preset": "custom",
		"retention_days":   3651,
	})
	if err == nil {
		t.Fatal("expected retention days too large")
	}
}
