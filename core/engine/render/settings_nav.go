package render

import (
	"context"
	"strconv"
	"strings"

	"sumeru/core/engine/parser"
	"sumeru/core/orm"
)

// MenuGeneralSettingsXMLID is the settings hub sidebar entry (URL action → SettingsHubPath).
const MenuGeneralSettingsXMLID = "base.menu_general_settings"

const (
	SettingsHubPath    = "/web/settings"
	SettingsAccountPath = "/web/settings/account"
)

func settingsNavExcludedSection(title string) bool {
	return strings.EqualFold(strings.TrimSpace(title), "Personal")
}

func settingsNavExcludedLink(name, href string) bool {
	if strings.EqualFold(strings.TrimSpace(name), "Account security") {
		return true
	}
	return strings.Contains(strings.TrimSpace(href), SettingsAccountPath)
}

func menuIDForXMLID(ctx context.Context, xmlID string) string {
	xmlID = strings.TrimSpace(xmlID)
	if xmlID == "" {
		return ""
	}
	id, _, err := orm.ResolveXmlId(ctx, xmlID)
	if err != nil || id <= 0 {
		return ""
	}
	return strconv.Itoa(id)
}

func isSettingsHubNavMenu(ctx context.Context, menuID string) bool {
	menuID = strings.TrimSpace(menuID)
	if menuID == "" {
		return false
	}
	return menuID == menuIDForXMLID(ctx, MenuGeneralSettingsXMLID)
}

// MenuIDIsGeneralSettings reports whether menuID is base.menu_general_settings.
func MenuIDIsGeneralSettings(ctx context.Context, menuID string) bool {
	return isSettingsHubNavMenu(ctx, menuID)
}

func applySettingsNavHrefOverrides(ctx context.Context, menus []parser.MenuItem) {
	generalID := menuIDForXMLID(ctx, MenuGeneralSettingsXMLID)
	if generalID == "" {
		return
	}
	for i := range menus {
		if menus[i].ID == generalID {
			menus[i].Action = SettingsHubPath
			return
		}
	}
}

// ApplySettingsNavHrefOverridesForTest applies in-memory href overrides (tests without DB xml ids).
func ApplySettingsNavHrefOverridesForTest(menus []parser.MenuItem, generalSettingsMenuID string) {
	generalSettingsMenuID = strings.TrimSpace(generalSettingsMenuID)
	if generalSettingsMenuID == "" {
		return
	}
	for i := range menus {
		if menus[i].ID == generalSettingsMenuID {
			menus[i].Action = SettingsHubPath
			return
		}
	}
}
