import { SWC_API_BASE } from "../constants/routes.js";

type BusHandler = (payload: unknown) => void;

interface BusEventFrame {
  type?: string;
  id?: number;
  channel?: string;
  payload?: unknown;
}

function parseBusEventFrame(raw: unknown): BusEventFrame | null {
  if (typeof raw !== "object" || raw === null) return null;
  return raw as BusEventFrame;
}

function lastEventStorageKey(channel: string): string {
  return `sum-bus-last:${channel}`;
}

function readLastEventId(channel: string): number {
  try {
    const v = sessionStorage.getItem(lastEventStorageKey(channel));
    if (!v) return 0;
    const n = Number(v);
    return Number.isFinite(n) && n > 0 ? n : 0;
  } catch {
    return 0;
  }
}

function writeLastEventId(channel: string, id: number): void {
  if (id <= 0) return;
  try {
    sessionStorage.setItem(lastEventStorageKey(channel), String(id));
  } catch {
    /* ignore */
  }
}

/** Client event bus with WebSocket live updates from /web/swc/bus. */
export class BusService {
  private readonly handlers = new Map<string, Set<BusHandler>>();
  private ws: WebSocket | null = null;
  private wsURL = `${SWC_API_BASE}/bus`;
  private reconnectAttempt = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private remoteChannels = new Set<string>();
  private intentionalClose = false;

  subscribe(channel: string, handler: BusHandler): () => void {
    if (!this.handlers.has(channel)) {
      this.handlers.set(channel, new Set());
    }
    this.handlers.get(channel)!.add(handler);
    return () => this.handlers.get(channel)?.delete(handler);
  }

  emit(channel: string, payload: unknown): void {
    for (const fn of this.handlers.get(channel) ?? []) {
      fn(payload);
    }
  }

  /** Subscribe on server when connected; always registers local handler via subscribe(). */
  watchChannel(channel: string, handler: BusHandler): () => void {
    const unsubLocal = this.subscribe(channel, handler);
    this.remoteChannels.add(channel);
    this.sendSubscribe([channel]);
    return () => {
      unsubLocal();
      this.remoteChannels.delete(channel);
      this.sendUnsubscribe([channel]);
    };
  }

  watchRecord(model: string, id: number, handler: BusHandler): () => void {
    if (!model || id <= 0) return () => undefined;
    const channel = `record/${model}/${id}`;
    return this.watchChannel(channel, handler);
  }

  connect(url = `${SWC_API_BASE}/bus`): void {
    this.wsURL = url;
    this.intentionalClose = false;
    this.openSocket();
  }

  disconnect(): void {
    this.intentionalClose = true;
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer);
      this.reconnectTimer = null;
    }
    this.ws?.close();
    this.ws = null;
  }

  private openSocket(): void {
    if (this.ws) return;
    try {
      const proto = window.location.protocol === "https:" ? "wss:" : "ws:";
      this.ws = new WebSocket(`${proto}//${window.location.host}${this.wsURL}`);
      this.ws.addEventListener("open", () => {
        this.reconnectAttempt = 0;
        if (this.remoteChannels.size > 0) {
          this.sendSubscribe([...this.remoteChannels]);
        }
      });
      this.ws.addEventListener("message", (ev) => {
        try {
          const parsed: unknown = JSON.parse(String(ev.data));
          const frame = parseBusEventFrame(parsed);
          if (!frame) return;
          if (frame.type === "event" && frame.channel) {
            if (typeof frame.id === "number" && frame.id > 0) {
              writeLastEventId(frame.channel, frame.id);
            }
            this.emit(frame.channel, frame.payload);
          }
        } catch (err) {
          console.warn("swc bus: malformed message", err);
        }
      });
      this.ws.addEventListener("close", () => {
        this.ws = null;
        if (!this.intentionalClose) {
          this.scheduleReconnect();
        }
      });
    } catch (err) {
      console.warn("swc bus: WebSocket unavailable; local-only bus", err);
      this.scheduleReconnect();
    }
  }

  private scheduleReconnect(): void {
    if (this.reconnectTimer || this.intentionalClose) return;
    const delay = Math.min(60_000, 1000 * 2 ** this.reconnectAttempt) * (0.8 + Math.random() * 0.4);
    this.reconnectAttempt += 1;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.openSocket();
    }, delay);
  }

  private sendFrame(obj: Record<string, unknown>): void {
    if (!this.ws || this.ws.readyState !== WebSocket.OPEN) return;
    this.ws.send(JSON.stringify(obj));
  }

  private sendSubscribe(channels: string[]): void {
    if (channels.length === 0) return;
    const frame: Record<string, unknown> = { type: "subscribe", channels };
    const maxLast = Math.max(0, ...channels.map((c) => readLastEventId(c)));
    if (maxLast > 0) {
      frame.last_event_id = maxLast;
    }
    this.sendFrame(frame);
  }

  private sendUnsubscribe(channels: string[]): void {
    if (channels.length === 0) return;
    this.sendFrame({ type: "unsubscribe", channels });
  }
}
