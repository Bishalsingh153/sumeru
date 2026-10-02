package web

import (
	"context"
	"net/http"
	"strings"

	"sumeru/core/orm"
	"sumeru/core/server/auth"
)

const settingsAccountTOTPQRRoute = "/web/settings/account/totp-qr"

func registerSettingsAccountTOTPRoutes() {
	registerSession(http.MethodGet, settingsAccountTOTPQRRoute, SettingsAccountTOTPQrHandler)
}

func SettingsAccountTOTPQrHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	ctx := r.Context()
	actor := orm.SecurityUID(ctx)
	secret, ok, err := orm.OwnTOTPEnrollmentPending(ctx, actor, actor)
	if err != nil || !ok {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	login := userLoginName(ctx, actor)
	uri := auth.TOTPOtpauthURI("Sumeru", login, secret)
	png, err := auth.TOTPQRPng(uri, 220)
	if err != nil {
		http.Error(w, "qr unavailable", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "no-store")
	_, _ = w.Write(png)
}

func userLoginName(ctx context.Context, userID int) string {
	rec, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"id": userID})
	if err != nil || rec == nil {
		return "user"
	}
	login := strings.TrimSpace(ormRowString(rec, "login"))
	if login == "" {
		return "user"
	}
	return login
}

func ormRowString(rec map[string]interface{}, key string) string {
	if v, ok := rec[key].(string); ok {
		return v
	}
	return ""
}

func handleSettingsAccountTOTP(w http.ResponseWriter, r *http.Request, ctx context.Context, actor int, action string) bool {
	switch strings.TrimSpace(action) {
	case "totp_enroll_start":
		if _, err := orm.BeginOwnTOTPEnrollment(ctx, actor, actor); err != nil {
			redirectWithWebMessage(w, r, settingsAccountRoute, "totp_enroll_failed")
			return true
		}
		redirectWithWebMessage(w, r, settingsAccountRoute, "totp_enroll_started")
		return true
	case "totp_enroll_confirm":
		code := strings.TrimSpace(r.PostForm.Get("totp_code"))
		if err := orm.ConfirmOwnTOTPEnrollment(ctx, actor, actor, code); err != nil {
			redirectWithWebMessage(w, r, settingsAccountRoute, "totp_invalid")
			return true
		}
		redirectWithWebMessage(w, r, settingsAccountRoute, "totp_enabled")
		return true
	case "totp_disable":
		code := strings.TrimSpace(r.PostForm.Get("totp_code"))
		if err := orm.DisableOwnTOTP(ctx, actor, actor, code); err != nil {
			redirectWithWebMessage(w, r, settingsAccountRoute, "totp_invalid")
			return true
		}
		redirectWithWebMessage(w, r, settingsAccountRoute, "totp_disabled")
		return true
	default:
		return false
	}
}

func loadSettingsAccountTOTPState(ctx context.Context, actor int) settingsAccountTOTPState {
	enabled, _ := orm.UserTOTPEnabled(ctx, actor)
	pending := false
	if !enabled {
		if _, ok, _ := orm.OwnTOTPEnrollmentPending(ctx, actor, actor); ok {
			pending = true
		}
	}
	return settingsAccountTOTPState{Enabled: enabled, EnrollPending: pending}
}

type settingsAccountTOTPState struct {
	Enabled       bool
	EnrollPending bool
}
