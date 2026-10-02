package web

import (
	"net/http"
	"strings"
	"time"

	"sumeru/core/orm"
)

const (
	swcNotificationsRoute    = "/web/swc/notifications"
	swcNotificationsReadRoute = "/web/swc/notifications/read"
)

func registerSwcNotificationRoutes() {
	registerSession(http.MethodGet, swcNotificationsRoute, SwcNotificationsHandler)
	registerSession(http.MethodPost, swcNotificationsReadRoute, SwcNotificationsMarkReadHandler)
}

type swcNotificationItem struct {
	ID          int    `json:"id"`
	MessageID   int    `json:"messageId"`
	Body        string `json:"body"`
	Author      string `json:"author"`
	RecordModel string `json:"recordModel"`
	RecordID    int    `json:"recordId"`
	IsRead      bool   `json:"isRead"`
	CreateDate  string `json:"createDate"`
}

// SwcNotificationsHandler GET /web/swc/notifications
func SwcNotificationsHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	uid := AuthenticatedUserID(r)
	ctx := r.Context()
	if err := orm.CheckModelAccess(ctx, uid, "mail.notification", "read"); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	unreadOnly := strings.TrimSpace(r.URL.Query().Get("unread")) == "1"
	domain := [][]interface{}{{"user_id", "=", uid}}
	if unreadOnly {
		domain = append(domain, []interface{}{"is_read", "=", false})
	}
	rows, err := orm.SearchPage(ctx, "mail.notification", domain, 50, 0, "create_date DESC")
	if err != nil {
		http.Error(w, "load failed", http.StatusInternalServerError)
		return
	}
	items := make([]swcNotificationItem, 0, len(rows))
	unread := 0
	for _, row := range rows {
		msgID := int(coerceFloat(row["message_id"]))
		body, author := "", ""
		if msgID > 0 {
			if msg, err := orm.SearchOne(ctx, "mail.message", map[string]interface{}{"id": msgID}); err == nil && msg != nil {
				body, _ = msg["body"].(string)
				author, _ = msg["author"].(string)
			}
		}
		isRead, _ := row["is_read"].(bool)
		if !isRead {
			unread++
		}
		items = append(items, swcNotificationItem{
			ID:          int(coerceFloat(row["id"])),
			MessageID:   msgID,
			Body:        body,
			Author:      author,
			RecordModel: notifStringField(row["record_model"]),
			RecordID:    int(coerceFloat(row["record_id"])),
			IsRead:      isRead,
			CreateDate:  notifStringField(row["create_date"]),
		})
	}
	writeJSONResponse(w, map[string]interface{}{
		"items":  items,
		"unread": unread,
	})
}

// SwcNotificationsMarkReadHandler POST /web/swc/notifications/read
func SwcNotificationsMarkReadHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !ParsePostForm(w, r) {
		return
	}
	uid := AuthenticatedUserID(r)
	ctx := r.Context()
	if err := orm.CheckModelAccess(ctx, uid, "mail.notification", "write"); err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}
	markAll := strings.TrimSpace(r.PostFormValue("all")) == "1"
	now := time.Now().UTC().Format(time.RFC3339)
	vals := map[string]interface{}{"is_read": true, "read_date": now}
	if markAll {
		rows, _ := orm.Search(ctx, "mail.notification", [][]interface{}{
			{"user_id", "=", uid},
			{"is_read", "=", false},
		})
		for _, row := range rows {
			id := int(coerceFloat(row["id"]))
			_, _ = orm.Update(ctx, "mail.notification", [][]interface{}{{"id", "=", id}}, vals)
		}
		writeJSONResponse(w, map[string]interface{}{"ok": true})
		return
	}
	id := int(coerceFloat(r.PostFormValue("id")))
	if id <= 0 {
		http.Error(w, "id required", http.StatusBadRequest)
		return
	}
	rec, err := orm.SearchOne(ctx, "mail.notification", map[string]interface{}{"id": id})
	if err != nil || rec == nil || int(coerceFloat(rec["user_id"])) != uid {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	_, err = orm.Update(ctx, "mail.notification", [][]interface{}{{"id", "=", id}}, vals)
	if err != nil {
		http.Error(w, "update failed", http.StatusInternalServerError)
		return
	}
	writeJSONResponse(w, map[string]interface{}{"ok": true})
}

func coerceFloat(v interface{}) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	case int64:
		return float64(t)
	default:
		return 0
	}
}

func notifStringField(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
