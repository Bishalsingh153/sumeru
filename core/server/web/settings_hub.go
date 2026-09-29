package web

import (
	"context"
	"net/http"

	"sumeru/core/engine/render"
	"sumeru/core/orm"
)

type settingsHubData struct {
	Categories   []render.SettingsHubCategory
	AppsListHref string
}

const settingsHubScriptURL = "/static/js/settings-hub.js"

// SettingsHubHandler renders the Settings overview at /web/settings.
func SettingsHubHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !requireSettingsUser(w, r) {
		return
	}

	ctx := r.Context()
	rootMenuIDStr, hubMenuIDStr, ok := resolveSettingsHubContext(w, r, ctx)
	if !ok {
		return
	}

	renderSettingsHubPage(w, r, settingsHubData{
		Categories:   render.BuildSettingsHubCategories(ctx, rootMenuIDStr),
		AppsListHref: appsRoute,
	}, hubMenuIDStr)
}

func requireSettingsUser(w http.ResponseWriter, r *http.Request) bool {
	ctx := r.Context()
	if orm.UserHasGroupXML(ctx, orm.SecurityUID(ctx), groupUserXML) {
		return true
	}
	http.Redirect(w, r, homeRoute, http.StatusFound)
	return false
}

func resolveSettingsRootMenuID(w http.ResponseWriter, r *http.Request, ctx context.Context) (menuIDStr string, ok bool) {
	_, menuIDStr = resolveMenuID(ctx, settingsRootMenuXMLID)
	if menuIDStr != "" {
		return menuIDStr, true
	}
	http.Redirect(w, r, appsRoute, http.StatusFound)
	return "", false
}

func resolveSettingsHubContext(w http.ResponseWriter, r *http.Request, ctx context.Context) (rootMenuIDStr, hubMenuIDStr string, ok bool) {
	_, rootMenuIDStr = resolveMenuID(ctx, settingsRootMenuXMLID)
	if rootMenuIDStr == "" {
		http.Redirect(w, r, appsRoute, http.StatusFound)
		return "", "", false
	}
	_, hubMenuIDStr = resolveMenuID(ctx, settingsHubMenuXMLID)
	if hubMenuIDStr == "" {
		hubMenuIDStr = rootMenuIDStr
	}
	return rootMenuIDStr, hubMenuIDStr, true
}

func renderSettingsHubPage(w http.ResponseWriter, r *http.Request, pageData settingsHubData, hubMenuIDStr string) {
	ctx := r.Context()
	renderShellPage(w, r, shellPageOpts{
		Route:         settingsRoute,
		InnerTemplate: settingsHubInnerTemplate,
		InnerData:     pageData,
		MenuIDStr:     hubMenuIDStr,
		Page:          buildSettingsHubPageData(ctx, hubMenuIDStr),
	})
}

func buildSettingsHubPageData(ctx context.Context, hubMenuIDStr string) render.PageData {
	return render.PageData{
		Title:             settingsHubPageTitle,
		SettingsNavActive: true,
		ActiveMenuID:      hubMenuIDStr,
		BreadcrumbItems:   render.BuildSettingsHubBreadcrumbs(ctx),
		ViewStylesheetURLs: []string{
			settingsHubStylesheetURL,
		},
		ExtraBodyClasses: settingsHubBodyClass,
		ExtraScriptURLs:  []string{settingsHubScriptURL},
	}
}
