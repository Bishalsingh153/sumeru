package web

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/lib/pq"
	"sumeru/core/applog"
	"sumeru/core/orm"
	"sumeru/core/queue"
)

const swcBusRoute = "/web/swc/bus"

var (
	swcBusUpgrader = websocket.Upgrader{
		CheckOrigin: checkSwcBusOrigin,
	}
	globalBusHub     *busHub
	globalBusHubOnce sync.Once
)

const (
	swcBusSendBuffer   = 32
	swcBusPingInterval = 45 * time.Second
	swcBusPongWait     = 60 * time.Second
)

func init() {
	orm.SetBusPublishHook(func(ctx context.Context, channel string, payload map[string]interface{}) {
		_ = PublishSwcBusEvent(ctx, channel, payload)
	})
}

func registerSwcBusRoute() {
	registerSession(http.MethodGet, swcBusRoute, SwcBusHandler)
}

// SwcBusHandler upgrades GET /web/swc/bus to a WebSocket for live outbox events.
func SwcBusHandler(w http.ResponseWriter, r *http.Request) {
	if !requireLogin(w, r) {
		return
	}
	if !websocket.IsWebSocketUpgrade(r) {
		http.NotFound(w, r)
		return
	}
	serveSwcBusWebSocket(w, r, AuthenticatedUserID(r))
}

// StartBusListener listens for PostgreSQL NOTIFY and dispatches bus events locally.
func StartBusListener(parent context.Context, connStr string) {
	if parent == nil || connStr == "" || orm.DB == nil {
		return
	}
	go func() {
		backoff := time.Second
		for {
			if parent.Err() != nil {
				return
			}
			if err := listenBusLoop(parent, connStr); err != nil {
				applog.WarnMsg(parent, "web", "bus_listen", "bus listener ended", err, nil)
				time.Sleep(backoff)
				if backoff < 30*time.Second {
					backoff *= 2
				}
				continue
			}
			return
		}
	}()
}

func listenBusLoop(parent context.Context, connStr string) error {
	listener := pq.NewListener(connStr, 10*time.Second, time.Minute, func(ev pq.ListenerEventType, err error) {
		if err != nil {
			applog.WarnMsg(parent, "web", "bus_listen", "pq listener event", err, map[string]interface{}{"event": int(ev)})
		}
	})
	defer func() { _ = listener.Close() }()
	if err := listener.Listen(orm.BusNotifyChannel()); err != nil {
		return err
	}
	for {
		select {
		case <-parent.Done():
			return parent.Err()
		case n, ok := <-listener.Notify:
			if !ok {
				return nil
			}
			if n != nil {
				dispatchBusEventByID(parent, parseBusEventIDNotify(n.Extra))
			}
		}
	}
}

func checkSwcBusOrigin(r *http.Request) bool {
	origin := strings.TrimSpace(r.Header.Get("Origin"))
	if origin == "" {
		return true
	}
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	reqHost := r.Host
	if h, _, err := net.SplitHostPort(reqHost); err == nil {
		reqHost = h
	}
	return strings.EqualFold(u.Hostname(), reqHost)
}

type swcBusClient struct {
	uid      int
	conn     *websocket.Conn
	send     chan []byte
	channels map[string]struct{}
	subMu    sync.Mutex
}

type busHub struct {
	mu      sync.RWMutex
	clients map[*swcBusClient]struct{}
}

func ensureBusHub() *busHub {
	globalBusHubOnce.Do(func() {
		globalBusHub = &busHub{clients: make(map[*swcBusClient]struct{})}
		queue.Subscribe("outbox", handleOutboxBusBridge)
	})
	return globalBusHub
}

func handleOutboxBusBridge(ctx context.Context, msg queue.Message) error {
	var envelope map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &envelope); err != nil {
		return nil
	}
	name, _ := envelope["name"].(string)
	if name == "" {
		return nil
	}
	inner, _ := envelope["payload"].(map[string]interface{})
	if inner == nil {
		inner = map[string]interface{}{}
	}
	model, _ := inner["model"].(string)
	var rid int
	switch v := inner["id"].(type) {
	case float64:
		rid = int(v)
	case int:
		rid = v
	}
	if model != "" && rid > 0 && strings.HasPrefix(name, "record.") {
		PublishRecordBusEvent(ctx, name, model, rid)
	}
	return nil
}

func (h *busHub) register(c *swcBusClient) {
	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()
}

func (h *busHub) unregister(c *swcBusClient) {
	h.mu.Lock()
	delete(h.clients, c)
	h.mu.Unlock()
}

func (h *busHub) publishChannel(channel string, msg []byte) {
	h.mu.RLock()
	targets := make([]*swcBusClient, 0)
	for c := range h.clients {
		c.subMu.Lock()
		_, ok := c.channels[channel]
		c.subMu.Unlock()
		if ok {
			targets = append(targets, c)
		}
	}
	h.mu.RUnlock()
	for _, c := range targets {
		select {
		case c.send <- msg:
		default:
		}
	}
}

func (c *swcBusClient) subscribeChannels(ctx context.Context, channels []string, afterID int64) {
	for _, ch := range channels {
		ch = strings.TrimSpace(ch)
		if ch == "" || !AuthorizeSwcBusChannel(ctx, c.uid, ch) {
			continue
		}
		c.subMu.Lock()
		if len(c.channels) >= maxBusSubscriptions {
			c.subMu.Unlock()
			break
		}
		c.channels[ch] = struct{}{}
		c.subMu.Unlock()
		if afterID > 0 {
			rows, err := orm.ListBusEventsAfter(ctx, []string{ch}, afterID, 500)
			if err != nil {
				continue
			}
			for _, row := range rows {
				frame, err := marshalBusEventFrame(row.ID, row.Channel, row.Payload)
				if err != nil {
					continue
				}
				select {
				case c.send <- frame:
				default:
				}
			}
		}
	}
}

func (c *swcBusClient) unsubscribeChannels(channels []string) {
	c.subMu.Lock()
	defer c.subMu.Unlock()
	for _, ch := range channels {
		ch = strings.TrimSpace(ch)
		delete(c.channels, ch)
	}
}

func (c *swcBusClient) writePump() {
	ticker := time.NewTicker(swcBusPingInterval)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *swcBusClient) readPump(h *busHub, baseCtx context.Context) {
	defer func() {
		h.unregister(c)
		close(c.send)
		c.conn.Close()
	}()
	_ = c.conn.SetReadDeadline(time.Now().Add(swcBusPongWait))
	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(swcBusPongWait))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		c.handleClientFrame(baseCtx, data)
	}
}

func (c *swcBusClient) handleClientFrame(ctx context.Context, data []byte) {
	var frame map[string]interface{}
	if err := json.Unmarshal(data, &frame); err != nil {
		return
	}
	typ, _ := frame["type"].(string)
	switch typ {
	case "subscribe":
		channels := stringSliceField(frame["channels"])
		afterID, _ := orm.CoerceInt64(frame["last_event_id"])
		c.subscribeChannels(ctx, channels, afterID)
	case "unsubscribe":
		c.unsubscribeChannels(stringSliceField(frame["channels"]))
	case "ping":
		if b, err := marshalBusFrame(map[string]interface{}{"type": "pong"}); err == nil {
			select {
			case c.send <- b:
			default:
			}
		}
	}
}

func stringSliceField(v interface{}) []string {
	switch t := v.(type) {
	case []interface{}:
		out := make([]string, 0, len(t))
		for _, item := range t {
			if s, ok := item.(string); ok && s != "" {
				out = append(out, s)
			}
		}
		return out
	case []string:
		return t
	default:
		return nil
	}
}

func marshalBusFrame(m map[string]interface{}) ([]byte, error) {
	return json.Marshal(m)
}

func serveSwcBusWebSocket(w http.ResponseWriter, r *http.Request, uid int) {
	hub := ensureBusHub()
	conn, err := swcBusUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	client := &swcBusClient{
		uid:      uid,
		conn:     conn,
		send:     make(chan []byte, swcBusSendBuffer),
		channels: make(map[string]struct{}),
	}
	hub.register(client)
	ctx := r.Context()
	client.subscribeChannels(ctx, []string{UserNotificationsChannel(uid)}, 0)
	go client.writePump()
	client.readPump(hub, ctx)
}
