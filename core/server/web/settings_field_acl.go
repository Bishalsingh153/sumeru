package web

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"sumeru/core/engine/render"
	"sumeru/core/orm"
)

const settingsFieldACLRoute = "/web/settings/field-acl"

type fieldACLMatrixCell struct {
	FieldName  string
	GroupID    int
	GroupLabel string
	DenyRead   bool
	DenyWrite  bool
	HasGroups  bool
}

type fieldACLMatrixRow struct {
	FieldName string
	HasGroups bool
	Cells     []fieldACLMatrixCell
}

type settingsFieldACLData struct {
	CSRFToken string
	Model     string
	Groups    []fieldACLGroupCol
	Rows      []fieldACLMatrixRow
	Flash     render.FlashMessage
}

type fieldACLGroupCol struct {
	ID    int
	Label string
}

func registerSettingsFieldACLRoutes() {
	registerSession(http.MethodGet, settingsFieldACLRoute, SettingsFieldACLGetHandler)
	registerSession(http.MethodPost, settingsFieldACLRoute, SettingsFieldACLPostHandler)
}

func SettingsFieldACLGetHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !requireSystemAdmin(w, r, true) {
		return
	}
	ctx := r.Context()
	menuIDStr, ok := resolveSettingsRootMenuID(w, r, ctx)
	if !ok {
		return
	}
	model := strings.TrimSpace(r.URL.Query().Get("model"))
	flash, _ := flashFromQueryMessage(r.URL.Query().Get("msg"))
	page, err := buildFieldACLMatrixPage(ctx, model)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	page.CSRFToken = CSRFTokenForRequest(r)
	page.Flash = flash
	renderSettingsFieldACLPage(w, r, page, menuIDStr)
}

func SettingsFieldACLPostHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !requireSystemAdmin(w, r, true) {
		return
	}
	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	if !ValidateCSRF(r) {
		http.Error(w, "invalid csrf", http.StatusForbidden)
		return
	}
	model := strings.TrimSpace(r.PostForm.Get("model"))
	if model == "" || orm.RegistryModel(model) == nil {
		http.Error(w, "model required", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	err := orm.WithElevated(ctx, "settings.field_acl_matrix", func(elevated context.Context) error {
		return applyFieldACLMatrixPost(elevated, model, r.PostForm)
	})
	if err != nil {
		redirectWithWebMessage(w, r, settingsFieldACLRoute+"?model="+model, "field_acl_failed")
		return
	}
	redirectWithWebMessage(w, r, settingsFieldACLRoute+"?model="+model, "field_acl_saved")
}

func applyFieldACLMatrixPost(ctx context.Context, model string, form map[string][]string) error {
	groups, err := loadFieldACLGroups(ctx)
	if err != nil {
		return err
	}
	groupIDs := []int{0}
	for _, g := range groups {
		groupIDs = append(groupIDs, g.ID)
	}
	fields := listFieldACLFields(model)
	existing, err := loadFieldACLExisting(ctx, model)
	if err != nil {
		return err
	}
	for _, field := range fields {
		for _, gid := range groupIDs {
			key := matrixCellKey(gid, field)
			denyRead := formHas(form, "dr_"+key)
			denyWrite := formHas(form, "dw_"+key)
			rowName := fieldACLMatrixRowName(model, field, gid)
			if !denyRead && !denyWrite {
				if id, ok := existing[rowName]; ok && id > 0 {
					if err := orm.Unlink(ctx, "sys.field.access", id); err != nil {
						return err
					}
				}
				continue
			}
			values := map[string]interface{}{
				"name":       rowName,
				"model":      model,
				"field_name": field,
				"perm_read":  !denyRead,
				"perm_write": !denyWrite,
			}
			if gid > 0 {
				values["group_id"] = gid
			}
			if _, err := orm.Upsert(ctx, orm.RegistryModel("sys.field.access"), values, "name"); err != nil {
				return err
			}
		}
	}
	return nil
}

func formHas(form map[string][]string, key string) bool {
	vals, ok := form[key]
	return ok && len(vals) > 0 && strings.TrimSpace(vals[0]) != ""
}

func matrixCellKey(groupID int, field string) string {
	return fmt.Sprintf("%d_%s", groupID, field)
}

func fieldACLMatrixRowName(model, field string, groupID int) string {
	if groupID <= 0 {
		return "matrix." + model + "." + field + ".global"
	}
	return fmt.Sprintf("matrix.%s.%s.g%d", model, field, groupID)
}

func buildFieldACLMatrixPage(ctx context.Context, model string) (settingsFieldACLData, error) {
	out := settingsFieldACLData{Model: model}
	if model == "" {
		return out, nil
	}
	if orm.RegistryModel(model) == nil {
		return out, fmt.Errorf("unknown model %q", model)
	}
	groups, err := loadFieldACLGroups(ctx)
	if err != nil {
		return out, err
	}
	out.Groups = groups
	existing, err := loadFieldACLExistingByCell(ctx, model)
	if err != nil {
		return out, err
	}
	for _, field := range listFieldACLFields(model) {
		row := fieldACLMatrixRow{FieldName: field}
		for _, g := range append([]fieldACLGroupCol{{ID: 0, Label: "Global"}}, groups...) {
			key := matrixCellKey(g.ID, field)
			cell := fieldACLMatrixCell{
				FieldName:  field,
				GroupID:    g.ID,
				GroupLabel: g.Label,
			}
			if st, ok := existing[key]; ok {
				cell.DenyRead = !st.read
				cell.DenyWrite = !st.write
			}
			if fieldHasGroupsAttr(model, field) {
				cell.HasGroups = true
				row.HasGroups = true
			}
			row.Cells = append(row.Cells, cell)
		}
		out.Rows = append(out.Rows, row)
	}
	return out, nil
}

type fieldACLPerm struct {
	read  bool
	write bool
}

func loadFieldACLExistingByCell(ctx context.Context, model string) (map[string]fieldACLPerm, error) {
	rows, err := orm.Search(ctx, "sys.field.access", [][]interface{}{{"model", "=", model}})
	if err != nil {
		return nil, err
	}
	out := map[string]fieldACLPerm{}
	for _, row := range rows {
		field := stringField(row["field_name"])
		if field == "" {
			continue
		}
		gid := intField(row["group_id"])
		key := matrixCellKey(gid, field)
		out[key] = fieldACLPerm{
			read:  boolField(row["perm_read"], true),
			write: boolField(row["perm_write"], true),
		}
	}
	return out, nil
}

func loadFieldACLExisting(ctx context.Context, model string) (map[string]int, error) {
	rows, err := orm.Search(ctx, "sys.field.access", [][]interface{}{{"model", "=", model}})
	if err != nil {
		return nil, err
	}
	out := map[string]int{}
	for _, row := range rows {
		name := stringField(row["name"])
		id := intField(row["id"])
		if name != "" && id > 0 {
			out[name] = id
		}
	}
	return out, nil
}

func loadFieldACLGroups(ctx context.Context) ([]fieldACLGroupCol, error) {
	rows, err := orm.Search(ctx, "core.group", nil)
	if err != nil {
		return nil, err
	}
	out := make([]fieldACLGroupCol, 0, len(rows))
	for _, row := range rows {
		id := intField(row["id"])
		if id <= 0 {
			continue
		}
		out = append(out, fieldACLGroupCol{ID: id, Label: stringField(row["name"])})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Label < out[j].Label })
	return out, nil
}

func listFieldACLFields(model string) []string {
	inst := orm.RegistryModel(model)
	if inst == nil {
		return nil
	}
	skip := map[string]bool{"id": true, "create_uid": true, "write_uid": true, "create_date": true, "write_date": true}
	var names []string
	for _, f := range inst.Fields() {
		n := strings.TrimSpace(f.Name)
		if n == "" || skip[n] {
			continue
		}
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}

func fieldHasGroupsAttr(model, field string) bool {
	inst := orm.RegistryModel(model)
	if inst == nil {
		return false
	}
	for _, f := range inst.Fields() {
		if strings.TrimSpace(f.Name) == field {
			return strings.TrimSpace(f.Groups) != ""
		}
	}
	return false
}

func stringField(v interface{}) string {
	if v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return strings.TrimSpace(t)
	default:
		return strings.TrimSpace(fmt.Sprint(t))
	}
}

func intField(v interface{}) int {
	switch t := v.(type) {
	case int:
		return t
	case int64:
		return int(t)
	case float64:
		return int(t)
	case string:
		n, _ := strconv.Atoi(strings.TrimSpace(t))
		return n
	default:
		return 0
	}
}

func boolField(v interface{}, defaultVal bool) bool {
	switch t := v.(type) {
	case bool:
		return t
	case string:
		return strings.EqualFold(strings.TrimSpace(t), "true") || strings.TrimSpace(t) == "1"
	default:
		return defaultVal
	}
}

func renderSettingsFieldACLPage(w http.ResponseWriter, r *http.Request, pageData settingsFieldACLData, menuIDStr string) {
	ctx := r.Context()
	renderShellPage(w, r, shellPageOpts{
		Route:         settingsFieldACLRoute,
		InnerTemplate: settingsFieldACLInnerTemplate,
		InnerData:     pageData,
		MenuIDStr:     menuIDStr,
		Page:          buildSettingsFieldACLPageData(ctx, menuIDStr, pageData),
	})
}

func buildSettingsFieldACLPageData(ctx context.Context, menuIDStr string, page settingsFieldACLData) render.PageData {
	crumbs := render.BuildSettingsHubBreadcrumbs(ctx)
	crumbs = append(crumbs, render.BreadcrumbItem{Label: "Field access matrix"})
	pd := render.PageData{
		Title:             "Field access matrix",
		SettingsNavActive: true,
		ActiveMenuID:      menuIDStr,
		BreadcrumbItems:   crumbs,
		ViewStylesheetURLs:   []string{settingsHubStylesheetURL},
		ExtraBodyClasses:     settingsHubBodyClass,
	}
	if page.Flash.Body != "" || page.Flash.Title != "" {
		pd.FlashMessages = []render.FlashMessage{page.Flash}
	}
	return pd
}
