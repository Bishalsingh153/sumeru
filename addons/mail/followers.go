package mail

import (
	"context"
	"fmt"
	"html"
	"regexp"
	"strings"

	pkgmail "sumeru/core/mail"
	"sumeru/core/orm"
	"time"
)

var mentionPattern = regexp.MustCompile(`@([a-zA-Z0-9._-]+)`)

// SubscribeUserFollower ensures user_id follows the record (idempotent).
func SubscribeUserFollower(ctx context.Context, model string, resID int64, userID int) error {
	if userID <= 0 || resID <= 0 || strings.TrimSpace(model) == "" {
		return nil
	}
	if _, ok := orm.Registry["mail.follower"]; !ok {
		return nil
	}
	existing, err := orm.SearchOne(ctx, "mail.follower", map[string]interface{}{
		"res_model": model,
		"res_id":    int(resID),
		"user_id":   userID,
	})
	if err == nil && existing != nil {
		return nil
	}
	inst := orm.Registry["mail.follower"]
	_, err = orm.Create(ctx, inst, map[string]interface{}{
		"res_model": model,
		"res_id":    int(resID),
		"user_id":   userID,
		"active":    true,
	})
	return err
}

// NotifyMessageFollowers creates inbox notifications and optional email for a posted message.
func NotifyMessageFollowers(ctx context.Context, messageID int, authorUserID int) error {
	if messageID <= 0 || orm.DB == nil {
		return nil
	}
	msg, err := orm.SearchOne(ctx, "mail.message", map[string]interface{}{"id": messageID})
	if err != nil || msg == nil {
		return err
	}
	model := rowString(msg, "model")
	resID := int(rowInt64(msg, "core_id"))
	subtype := rowString(msg, "subtype")
	body := rowString(msg, "body")
	if model == "" || resID <= 0 {
		return nil
	}
	recipients := map[int]struct{}{}
	followers, _ := orm.Search(ctx, "mail.follower", [][]interface{}{
		{"res_model", "=", model},
		{"res_id", "=", resID},
		{"active", "=", true},
	})
	for _, f := range followers {
		uid := int(rowInt64(f, "user_id"))
		if uid <= 0 || uid == authorUserID {
			continue
		}
		recipients[uid] = struct{}{}
	}
	for _, login := range mentionPattern.FindAllStringSubmatch(body, -1) {
		if len(login) < 2 {
			continue
		}
		uid := userIDByLogin(ctx, login[1])
		if uid > 0 && uid != authorUserID {
			recipients[uid] = struct{}{}
			_ = SubscribeUserFollower(ctx, model, int64(resID), uid)
		}
	}
	authorName := rowString(msg, "author")
	subject := fmt.Sprintf("[%s] %s", model, truncate(body, 80))
	for uid := range recipients {
		if err := orm.CheckModelAccess(ctx, uid, model, "read"); err != nil {
			continue
		}
		notifID, err := createNotification(ctx, uid, messageID, model, resID)
		if err != nil {
			continue
		}
		orm.PublishBusChannel(ctx, orm.UserNotificationsBusChannel(uid), map[string]interface{}{
			"notification_id": notifID,
			"message_id":      messageID,
			"model":           model,
			"id":              resID,
		})
		if shouldEmailUser(ctx, uid) {
			email := userEmail(ctx, uid)
			if email != "" {
				htmlBody := renderNotificationHTML(authorName, body, model, resID)
				pkgmail.EnqueueHTML(ctx, email, subject, htmlBody, stripHTML(body))
			}
		}
	}
	_ = subtype
	return nil
}

func createNotification(ctx context.Context, uid, messageID int, model string, resID int) (int, error) {
	inst := orm.Registry["mail.notification"]
	id, err := orm.Create(ctx, inst, map[string]interface{}{
		"user_id":      uid,
		"message_id":   messageID,
		"is_read":      false,
		"record_model": model,
		"record_id":    resID,
		"create_date":  time.Now().UTC().Format(time.RFC3339),
	})
	return id, err
}

func userIDByLogin(ctx context.Context, login string) int {
	login = strings.TrimSpace(login)
	if login == "" || orm.DB == nil {
		return 0
	}
	rec, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"login": login})
	if err != nil || rec == nil {
		return 0
	}
	return int(rowInt64(rec, "id"))
}

func userEmail(ctx context.Context, uid int) string {
	rec, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"id": uid})
	if err != nil || rec == nil {
		return ""
	}
	return rowString(rec, "email")
}

func shouldEmailUser(ctx context.Context, uid int) bool {
	rec, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"id": uid})
	if err != nil || rec == nil {
		return false
	}
	if v, ok := rec["notify_email"].(bool); ok {
		return v
	}
	return true
}

func renderNotificationHTML(author, body, model string, resID int) string {
	esc := html.EscapeString(body)
	authorEsc := html.EscapeString(author)
	return "<html><body><p><strong>" + authorEsc + "</strong> on " + html.EscapeString(model) +
		" #" + fmt.Sprint(resID) + ":</p><p>" + strings.ReplaceAll(esc, "\n", "<br/>") + "</p></body></html>"
}

func stripHTML(s string) string {
	return strings.TrimSpace(s)
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
