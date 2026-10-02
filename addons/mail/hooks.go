package mail

import (
	"context"

	"sumeru/core/event"
	"sumeru/core/orm"
)

func init() {
	event.Subscribe(eventRecordCreated, onMailThreadRecordCreated)
	event.Subscribe(eventRecordUpdated, onMailThreadRecordUpdated)
}

const (
	eventRecordCreated = "record.created"
	eventRecordUpdated = "record.updated"
)

func onMailThreadRecordCreated(ctx context.Context, ev event.Event) error {
	model, _ := ev.Payload["model"].(string)
	id, _ := payloadInt(ev.Payload["id"])
	if model == "" || id <= 0 || !orm.ModelHasMailThread(model) {
		return nil
	}
	uid := ev.Actor
	if uid <= 0 {
		uid = orm.SecurityUID(ctx)
	}
	return SubscribeUserFollower(ctx, model, int64(id), uid)
}

func onMailThreadRecordUpdated(ctx context.Context, ev event.Event) error {
	model, _ := ev.Payload["model"].(string)
	id, _ := payloadInt(ev.Payload["id"])
	if model == "" || id <= 0 || !orm.ModelHasMailThread(model) {
		return nil
	}
	rec, err := orm.SearchOne(ctx, model, map[string]interface{}{"id": id})
	if err != nil || rec == nil {
		return nil
	}
	assignee, ok := payloadInt(rec["user_id"])
	if !ok || assignee <= 0 {
		return nil
	}
	return SubscribeUserFollower(ctx, model, int64(id), assignee)
}

func payloadInt(v interface{}) (int, bool) {
	switch t := v.(type) {
	case int:
	 return t, true
	case int64:
		return int(t), true
	case float64:
		return int(t), true
	default:
		return 0, false
	}
}
