package web

import (
	"net/http"
	"strconv"
	"strings"

	"sumeru/core/engine/render"
)

// ImportWizardPageHandler GET /web/import?batch= — shell page for SWC import wizard.
func ImportWizardPageHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	batchID := strings.TrimSpace(r.URL.Query().Get("batch"))
	if batchID == "" {
		http.Error(w, "batch required", http.StatusBadRequest)
		return
	}
	if _, err := strconv.Atoi(batchID); err != nil {
		http.Error(w, "invalid batch", http.StatusBadRequest)
		return
	}
	renderShellPage(w, r, shellPageOpts{
		Route:         importWizardRoute,
		InnerTemplate: render.TemplatePagesImportWizard,
		InnerData: map[string]interface{}{
			"BatchID":   batchID,
			"CSRFToken": CSRFTokenForRequest(r),
		},
		Page: render.PageData{Title: "Import data"},
	})
}
