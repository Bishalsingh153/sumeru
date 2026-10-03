package web_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"sumeru/core/server/web"
)

func TestSwcDirectPostRequiresLogin(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/web/swc/direct/post", nil)
	web.SwcDirectPostHandlerForTest(rec, req)
	if rec.Code != http.StatusFound && rec.Code != http.StatusForbidden && rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestSwcDirectUsersRequiresLogin(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/direct/users", nil)
	web.SwcDirectUsersHandlerForTest(rec, req)
	if rec.Code != http.StatusFound && rec.Code != http.StatusForbidden {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestSwcDirectConversationsWithSessionNoDB(t *testing.T) {
	web.SetTestSessionUserIDForTest(1)
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/direct/conversations", nil)
	web.SwcDirectConversationsHandlerForTest(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d", rec.Code)
	}
}

func TestSwcDirectUsersWithSessionNoDB(t *testing.T) {
	web.SetTestSessionUserIDForTest(1)
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/direct/users?q=admin", nil)
	web.SwcDirectUsersHandlerForTest(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body.String())
	}
}

func TestSwcDirectThreadRequiresUserId(t *testing.T) {
	web.SetTestSessionUserIDForTest(1)
	t.Cleanup(web.ResetTestSessionUserIDForTest)
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/web/swc/direct/thread", nil)
	web.SwcDirectThreadHandlerForTest(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatal("expected error without userId")
	}
}
