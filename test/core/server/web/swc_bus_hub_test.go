package web_test

import (
	"encoding/json"
	"testing"

	"sumeru/core/server/web"
)

func TestSwcBusHubPublishChannelTargetsSubscribers(t *testing.T) {
	h := web.NewBusHubForTest()
	userA := web.NewSwcBusClientForTest(1, 1)
	userB := web.NewSwcBusClientForTest(2, 1)
	h.Register(userA)
	h.Register(userB)
	userA.SubscribeChannel("record/crm.lead/1")

	msg := []byte(`{"type":"event","id":1,"channel":"record/crm.lead/1","payload":{"model":"crm.lead","id":1}}`)
	h.PublishChannel("record/crm.lead/1", msg)

	select {
	case got := <-userA.Recv():
		if string(got) != string(msg) {
			t.Fatalf("user A: got %q", got)
		}
	default:
		t.Fatal("expected message for subscribed user A")
	}
	select {
	case <-userB.Recv():
		t.Fatal("user B should not receive without subscription")
	default:
	}
}

func TestSwcBusHubEventFrameShape(t *testing.T) {
	h := web.NewBusHubForTest()
	client := web.NewSwcBusClientForTest(5, 1)
	h.Register(client)
	client.SubscribeChannel("model/core.partner")

	out, _ := json.Marshal(map[string]interface{}{
		"type":    "event",
		"id":      2,
		"channel": "model/core.partner",
		"payload": map[string]interface{}{"model": "core.partner", "id": 3},
	})
	h.PublishChannel("model/core.partner", out)

	select {
	case got := <-client.Recv():
		var parsed map[string]interface{}
		if err := json.Unmarshal(got, &parsed); err != nil {
			t.Fatal(err)
		}
		if parsed["channel"] != "model/core.partner" {
			t.Fatalf("channel: %v", parsed["channel"])
		}
	default:
		t.Fatal("expected bus message")
	}
}
