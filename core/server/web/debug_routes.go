package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"sumeru/core/orm"
)

const debugAccessRoute = "/web/debug/access"
const debugMetadataRoute = "/web/debug/metadata"

func registerDebugRoutes() {
	registerSession(http.MethodGet, debugAccessRoute, DebugAccessHandler)
	registerSession(http.MethodGet, debugMetadataRoute, DebugMetadataHandler)
}

// DebugAccessHandler GET /web/debug/access?model=
func DebugAccessHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !validateSessionCSRF(w, r) {
		return
	}
	if !requireSystemAdmin(w, r, false) {
		return
	}
	model := strings.TrimSpace(r.URL.Query().Get("model"))
	trace, err := orm.BuildDebugAccessTrace(r.Context(), orm.UIDFromContext(r.Context()), model)
	if err != nil {
		if _, ok := err.(*orm.AccessDeniedError); ok {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(trace)
}

// DebugMetadataHandler GET /web/debug/metadata?model=&record_id=
func DebugMetadataHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !validateSessionCSRF(w, r) {
		return
	}
	if !requireSystemAdmin(w, r, false) {
		return
	}
	model := strings.TrimSpace(r.URL.Query().Get("model"))
	recordID, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("record_id")))
	meta, err := orm.BuildDebugRecordMetadata(r.Context(), orm.UIDFromContext(r.Context()), model, recordID)
	if err != nil {
		if _, ok := err.(*orm.AccessDeniedError); ok {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(meta)
}
