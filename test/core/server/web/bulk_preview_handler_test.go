package web_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"sumeru/core/server/web"
)

func TestSwcImportHandlerMissingBatchRecord(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodGet, "/web/swc/import?batch=1", nil)
	rec := httptest.NewRecorder()
	web.SwcImportHandlerForTest(rec, req)
	if rec.Code != http.StatusNotFound && rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestSwcImportHandlerRequiresBatch(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodGet, "/web/swc/import", nil)
	rec := httptest.NewRecorder()
	web.SwcImportHandlerForTest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestImportWizardPageRequiresBatch(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodGet, "/web/import", nil)
	rec := httptest.NewRecorder()
	web.ImportWizardPageHandlerForTest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestExportAsyncHandlerRequiresModel(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodPost, "/web/export/async", bytes.NewReader([]byte("csrf_token=x")))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rec := httptest.NewRecorder()
	web.ExportAsyncHandlerForTest(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatal("expected failure without model")
	}
}

func TestBulkDryRunHandlerRejectsBadBody(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodPost, "/web/bulk/dry-run", bytes.NewReader([]byte(`{"batch_id":0}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	web.BulkDryRunHandlerForTest(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatal("expected failure")
	}
}

func TestBulkSaveImportTemplateHandlerRejectsBadBody(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodPost, "/web/bulk/save-template", bytes.NewReader([]byte(`{"batch_id":1}`)))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	web.BulkSaveImportTemplateHandlerForTest(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatal("expected failure without name")
	}
}

func TestBulkPreviewHandlerRejectsBadBody(t *testing.T) {
	web.ResetTestSessionUserIDForTest()
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	web.SetTestSessionUserIDForTest(1)
	req := httptest.NewRequest(http.MethodPost, "/web/bulk/preview", bytes.NewReader([]byte("{}")))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	web.BulkPreviewHandlerForTest(rec, req)
	if rec.Code != http.StatusBadRequest && rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d", rec.Code)
	}
}
