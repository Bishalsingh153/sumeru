package web

import (
	"encoding/json"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"sumeru/addons/mail"
	"sumeru/core/orm"
)

const chatterUploadRoute = "/web/chatter/upload"

// ChatterUploadHandler POST multipart: model, res_id, file → JSON attachment metadata.
func ChatterUploadHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLoginMultipartPost(w, r, maxImportBodyBytes) {
		return
	}
	if !mail.CompanyChatterEnabled(r.Context()) {
		http.Error(w, "chatter disabled", http.StatusForbidden)
		return
	}
	modelName := strings.TrimSpace(r.FormValue(recordModelField))
	resIDRaw := strings.TrimSpace(r.FormValue(chatterRecordIDField))
	resID, _ := strconv.ParseInt(resIDRaw, 10, 64)
	if modelName == "" || resID <= 0 {
		http.Error(w, "model and res_id required", http.StatusBadRequest)
		return
	}
	if !requireModelAccess(w, r, modelName, "write") {
		return
	}
	if !chatterTargetRecordExists(r, modelName, resID) {
		http.Error(w, "record not found", http.StatusNotFound)
		return
	}
	upload, header, err := r.FormFile(importFileField)
	if err != nil {
		http.Error(w, "file required", http.StatusBadRequest)
		return
	}
	defer upload.Close()
	content, err := io.ReadAll(upload)
	if err != nil || len(content) == 0 {
		http.Error(w, "empty file", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(header.Filename)
	if name == "" {
		name = "upload"
	}
	name = filepath.Base(name)
	attID, err := orm.CreateBinaryAttachment(r.Context(), orm.CreateBinaryAttachmentInput{
		Name:     name,
		ResModel: modelName,
		ResID:    int(resID),
		Data:     content,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   attID,
		"name": name,
		"url":  contentURL(attID, true),
	})
}

func contentURL(attachmentID int, inline bool) string {
	u := "/web/content/" + strconv.Itoa(attachmentID)
	if inline {
		u += "?download=0"
	}
	return u
}
