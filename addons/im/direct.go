package im

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"sumeru/core/orm"
)

const MessageModel = "im.message"

// SearchInternalUsers finds active internal users by name or login (excluding uid).
func SearchInternalUsers(ctx context.Context, uid int, query string, limit int) ([]map[string]interface{}, error) {
	if orm.DB == nil {
		return nil, nil
	}
	if err := orm.CheckModelAccess(ctx, uid, "core.user", "read"); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 40 {
		limit = 20
	}
	q := strings.TrimSpace(query)
	domain := [][]interface{}{
		{"active", "=", true},
		{"user_type", "=", "internal"},
		{"id", "!=", uid},
	}
	if q != "" {
		domain = append(domain, []interface{}{"|"})
		domain = append(domain, []interface{}{"name", "ilike", q})
		domain = append(domain, []interface{}{"login", "ilike", q})
	}
	return orm.SearchPage(ctx, "core.user", domain, limit, 0, "name ASC")
}

// ListDirectConversations returns recent peers for uid with last message preview.
func ListDirectConversations(ctx context.Context, uid int, limit int) ([]map[string]interface{}, error) {
	if orm.DB == nil {
		return nil, nil
	}
	if err := orm.CheckModelAccess(ctx, uid, MessageModel, "read"); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 50 {
		limit = 30
	}
	domain := [][]interface{}{
		{"|"},
		{"sender_id", "=", uid},
		{"recipient_id", "=", uid},
	}
	rows, err := orm.SearchPage(ctx, MessageModel, domain, 200, 0, "create_date DESC")
	if err != nil {
		return nil, err
	}
	type conv struct {
		preview string
		when    string
	}
	seen := map[int]conv{}
	order := []int{}
	for _, row := range rows {
		sid, _ := orm.CoerceInt64(row["sender_id"])
		rid, _ := orm.CoerceInt64(row["recipient_id"])
		var peer int
		switch uid {
		case int(sid):
			peer = int(rid)
		case int(rid):
			peer = int(sid)
		default:
			continue
		}
		if peer <= 0 || peer == uid {
			continue
		}
		if _, ok := seen[peer]; ok {
			continue
		}
		preview := strings.TrimSpace(orm.AsString(row["body"]))
		if len(preview) > 80 {
			preview = preview[:77] + "…"
		}
		seen[peer] = conv{
			preview: preview,
			when:    strings.TrimSpace(orm.AsString(row["create_date"])),
		}
		order = append(order, peer)
		if len(order) >= limit {
			break
		}
	}
	out := make([]map[string]interface{}, 0, len(order))
	for _, peerID := range order {
		c := seen[peerID]
		name := UserDisplayName(ctx, peerID)
		out = append(out, map[string]interface{}{
			"userId":      peerID,
			"name":        name,
			"preview":     c.preview,
			"lastMessage": c.when,
		})
	}
	return out, nil
}

// ListDirectThread returns messages between uid and peerID oldest first.
func ListDirectThread(ctx context.Context, uid, peerID int) ([]map[string]interface{}, error) {
	if peerID <= 0 || peerID == uid {
		return nil, fmt.Errorf("invalid peer")
	}
	if orm.DB == nil {
		return nil, nil
	}
	if err := orm.CheckModelAccess(ctx, uid, MessageModel, "read"); err != nil {
		return nil, err
	}
	if err := assertInternalUser(ctx, uid, peerID); err != nil {
		return nil, err
	}
	rowsA, err := orm.SearchPage(ctx, MessageModel, [][]interface{}{
		{"sender_id", "=", uid},
		{"recipient_id", "=", peerID},
	}, 120, 0, "create_date ASC")
	if err != nil {
		return nil, err
	}
	rowsB, err := orm.SearchPage(ctx, MessageModel, [][]interface{}{
		{"sender_id", "=", peerID},
		{"recipient_id", "=", uid},
	}, 120, 0, "create_date ASC")
	if err != nil {
		return nil, err
	}
	rows := append(rowsA, rowsB...)
	sortDirectMessages(rows)
	out := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		msgID, _ := orm.CoerceInt64(row["id"])
		sid, _ := orm.CoerceInt64(row["sender_id"])
		out = append(out, map[string]interface{}{
			"id":          int(msgID),
			"body":        strings.TrimSpace(orm.AsString(row["body"])),
			"createDate":  strings.TrimSpace(orm.AsString(row["create_date"])),
			"outgoing":    int(sid) == uid,
			"attachments": listMessageAttachments(ctx, uid, int(msgID)),
		})
	}
	return out, nil
}

func listMessageAttachments(ctx context.Context, uid, messageID int) []map[string]interface{} {
	if messageID <= 0 {
		return nil
	}
	if err := orm.CheckModelAccess(ctx, uid, "sys.attachment", "read"); err != nil {
		return nil
	}
	rows, err := orm.Search(ctx, "sys.attachment", [][]interface{}{
		{"model", "=", MessageModel},
		{"res_id", "=", messageID},
	})
	if err != nil {
		return nil
	}
	out := make([]map[string]interface{}, 0, len(rows))
	for _, row := range rows {
		id, _ := orm.CoerceInt64(row["id"])
		out = append(out, map[string]interface{}{
			"id":   int(id),
			"name": strings.TrimSpace(orm.AsString(row["name"])),
			"url":  fmt.Sprintf("/web/content/%d", id),
		})
	}
	return out
}

// PostDirectMessage sends an internal direct message from uid to recipientID.
func PostDirectMessage(ctx context.Context, uid, recipientID int, body string) (int, error) {
	body = strings.TrimSpace(body)
	if body == "" {
		return 0, fmt.Errorf("message required")
	}
	if err := orm.CheckModelAccess(ctx, uid, MessageModel, "create"); err != nil {
		return 0, err
	}
	if err := assertInternalUser(ctx, uid, recipientID); err != nil {
		return 0, err
	}
	inst, ok := orm.Registry[MessageModel]
	if !ok {
		return 0, fmt.Errorf("internal chat not available")
	}
	id, err := orm.Create(ctx, inst, map[string]interface{}{
		"sender_id":    uid,
		"recipient_id": recipientID,
		"body":         body,
		"create_date":  time.Now().UTC(),
	})
	if err != nil {
		return 0, err
	}
	return id, nil
}

func assertInternalUser(ctx context.Context, uid, peerID int) error {
	if err := orm.CheckModelAccess(ctx, uid, "core.user", "read"); err != nil {
		return err
	}
	u, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"id": peerID})
	if err != nil {
		return fmt.Errorf("user not found")
	}
	if !orm.AsBool(u["active"]) {
		return fmt.Errorf("user not available")
	}
	if strings.TrimSpace(orm.AsString(u["user_type"])) != "internal" {
		return fmt.Errorf("direct messages are for internal users only")
	}
	return nil
}

func sortDirectMessages(rows []map[string]interface{}) {
	sort.Slice(rows, func(i, j int) bool {
		a := strings.TrimSpace(orm.AsString(rows[i]["create_date"]))
		b := strings.TrimSpace(orm.AsString(rows[j]["create_date"]))
		return a < b
	})
}

// UserDisplayName returns a label for a user id.
func UserDisplayName(ctx context.Context, userID int) string {
	u, err := orm.SearchOne(ctx, "core.user", map[string]interface{}{"id": userID})
	if err != nil {
		return fmt.Sprintf("User #%d", userID)
	}
	if nm := strings.TrimSpace(orm.AsString(u["name"])); nm != "" {
		return nm
	}
	if lg := strings.TrimSpace(orm.AsString(u["login"])); lg != "" {
		return lg
	}
	return fmt.Sprintf("User #%d", userID)
}

// SortDirectMessagesForTest exposes sort for tests.
func SortDirectMessagesForTest(rows []map[string]interface{}) {
	sortDirectMessages(rows)
}
