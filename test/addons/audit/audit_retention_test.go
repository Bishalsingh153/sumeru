package audit_test

import (
	"context"
	"testing"

	"sumeru/addons/audit"
)

func TestRetentionOptionsEmptyPolicyDefaults(t *testing.T) {
	opts := audit.RetentionOptionsFromPolicyForTest(map[string]interface{}{}, false)
	if opts.RetentionDays != 90 || opts.BatchSize != 0 {
		t.Fatalf("opts: %+v", opts)
	}
}

func TestRetentionOptionsFromPolicyDefaults(t *testing.T) {
	opts := audit.RetentionOptionsFromPolicyForTest(map[string]interface{}{
		"hot_retention_days":   45,
		"delete_after_export":  true,
		"compress":             true,
	}, false)
	if opts.RetentionDays != 45 || !opts.DeleteAfterExport || !opts.Compress {
		t.Fatalf("opts: %+v", opts)
	}
}

func TestRetentionOptionsFromPolicyDryRun(t *testing.T) {
	opts := audit.RetentionOptionsFromPolicyForTest(map[string]interface{}{
		"hot_retention_days": 90,
		"dry_run_next":       true,
	}, false)
	if !opts.DryRun {
		t.Fatal("expected dry run")
	}
}

func TestAuditRetentionCronWithoutPolicy(t *testing.T) {
	if err := audit.RunAuditRetentionCronForTest(context.Background()); err != nil {
		t.Fatal(err)
	}
}

func TestRunAuditRetentionDryRun(t *testing.T) {
	result, err := audit.RunAuditRetentionForTest(context.Background(), audit.RetentionOptions{
		RetentionDays: 30,
		DryRun:        true,
		BatchSize:     100,
		MaxRowsPerRun: 1000,
	})
	if err != nil {
		t.Skip("dry run needs database:", err)
	}
	if result.Summary == "" {
		t.Fatal("expected summary")
	}
}

func TestRetentionOptionsLegalHold(t *testing.T) {
	opts := audit.RetentionOptionsFromPolicyForTest(map[string]interface{}{
		"hot_retention_days": 90,
		"legal_hold":         true,
	}, false)
	if !opts.LegalHold {
		t.Fatal("expected legal hold")
	}
}

func TestRetentionOptionsCompressDefault(t *testing.T) {
	opts := audit.RetentionOptionsFromPolicyForTest(map[string]interface{}{
		"hot_retention_days": 90,
	}, false)
	if !opts.Compress {
		t.Fatal("compress should default true")
	}
}
