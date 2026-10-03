package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sumeru/core/server/web"
)

func TestSwcActivityLogHandlerRequiresParams(t *testing.T) {
	web.SetTestSessionUserIDForTest(1)
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/activity-log", nil)
	web.SwcActivityLogHandlerForTest(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestSwcActivityLogHandlerRequiresLogin(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/activity-log?model=core.user&id=1", nil)
	web.SwcActivityLogHandlerForTest(rec, req)
	if rec.Code != http.StatusFound && rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d", rec.Code)
	}
}
