package audit_test

import (
	"testing"

	"sumeru/addons/audit"
)

func TestAuditRowFromMap(t *testing.T) {
	row := audit.AuditRowFromMapForTest(map[string]interface{}{
		"id": 9, "action": "write", "model": "core.user", "res_id": 3, "user_id": 1,
		"create_date": "2026-01-01T00:00:00Z", "detail": "ok",
	})
	if row.ID != 9 || row.Model != "core.user" || row.Detail != "ok" {
		t.Fatalf("row: %+v", row)
	}
}

func TestAuditPlainNDJSONMime(t *testing.T) {
	raw := []byte("{}\n")
	out, mime, err := audit.MaybeGzipAuditForTest(false, raw)
	if err != nil || mime != "application/x-ndjson" || string(out) != string(raw) {
		t.Fatalf("plain: mime=%s err=%v", mime, err)
	}
}

func TestAuditGzipExport(t *testing.T) {
	raw, err := audit.EncodeAuditJSONLForTest([]audit.AuditExportRowForTest{{ID: 1, Action: "read"}})
	if err != nil {
		t.Fatal(err)
	}
	gz, mime, err := audit.MaybeGzipAuditForTest(true, raw)
	if err != nil || mime != "application/gzip" || len(gz) == 0 {
		t.Fatalf("gzip: mime=%s err=%v", mime, err)
	}
	if audit.SHA256HexAuditForTest(gz) == "" {
		t.Fatal("sha256")
	}
}

func TestAuditBoundsAndIDs(t *testing.T) {
	rows := []audit.AuditExportRowForTest{
		{ID: 5, CreateDate: "2026-02-01T00:00:00Z"},
		{ID: 2, CreateDate: "2026-01-01T00:00:00Z"},
	}
	from, to := audit.AuditPeriodBoundsForTest(rows)
	if from != "2026-01-01T00:00:00Z" || to != "2026-02-01T00:00:00Z" {
		t.Fatalf("bounds: %s %s", from, to)
	}
	minID, maxID := audit.AuditMinMaxIDForTest(rows)
	if minID != 2 || maxID != 5 {
		t.Fatalf("ids: %d %d", minID, maxID)
	}
}

func TestAuditJSONLRoundTrip(t *testing.T) {
	rows := []audit.AuditExportRowForTest{
		{ID: 1, Action: "write", Model: "core.user", CreateDate: "2026-01-01T00:00:00Z"},
		{ID: 2, Action: "unlink", Model: "core.user", CreateDate: "2026-01-02T00:00:00Z"},
	}
	raw, err := audit.EncodeAuditJSONLForTest(rows)
	if err != nil {
		t.Fatal(err)
	}
	decoded, err := audit.DecodeAuditJSONLForTest(raw)
	if err != nil || len(decoded) != 2 || decoded[0].ID != 1 {
		t.Fatalf("decode: %+v err=%v", decoded, err)
	}
}
