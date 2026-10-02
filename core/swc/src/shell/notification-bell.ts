import type { BusService } from "../services/bus.js";
import type { HttpService } from "../services/http.js";
import { SWC_API_BASE } from "../constants/routes.js";

interface InboxPayload {
  items: Array<{
    id: number;
    body: string;
    author: string;
    recordModel: string;
    recordId: number;
    isRead: boolean;
  }>;
  unread: number;
}

export function initNotificationBell(http: HttpService, bus: BusService, userId: number): void {
  const host = document.querySelector(".sum-topbar-right");
  if (!host || userId <= 0) return;

  const btn = document.createElement("button");
  btn.type = "button";
  btn.className = "sum-icon-btn sum-notification-bell";
  btn.setAttribute("aria-label", "Notifications");
  btn.innerHTML = `<span class="sum-notification-bell-icon" aria-hidden="true">●</span><span class="sum-notification-badge"></span>`;

  const panel = document.createElement("div");
  panel.className = "sum-notification-panel";
  panel.hidden = true;

  host.insertBefore(btn, host.firstChild);
  host.insertBefore(panel, btn.nextSibling);

  const badge = btn.querySelector(".sum-notification-badge") as HTMLElement;
  const base = SWC_API_BASE;

  async function refresh(): Promise<void> {
    try {
      const data = await http.getJSON<InboxPayload>(`${base}/notifications?unread=0`);
      badge.textContent = data.unread > 0 ? String(data.unread) : "";
      badge.hidden = data.unread <= 0;
      panel.innerHTML = "";
      for (const item of data.items.slice(0, 20)) {
        const row = document.createElement("button");
        row.type = "button";
        row.className = "sum-notification-row" + (item.isRead ? "" : " sum-notification-row--unread");
        row.textContent = `${item.author}: ${item.body}`.slice(0, 120);
        row.addEventListener("click", () => {
          void http.postForm(`${base}/notifications/read`, { id: String(item.id) });
          if (item.recordModel && item.recordId > 0) {
            window.location.assign(`/web#action=&model=${encodeURIComponent(item.recordModel)}&id=${item.recordId}`);
          }
          void refresh();
        });
        panel.appendChild(row);
      }
      const markAll = document.createElement("button");
      markAll.type = "button";
      markAll.className = "sum-notification-mark-all";
      markAll.textContent = "Mark all read";
      markAll.addEventListener("click", () => {
        void http.postForm(`${base}/notifications/read`, { all: "1" }).then(refresh);
      });
      panel.appendChild(markAll);
    } catch {
      /* ignore */
    }
  }

  btn.addEventListener("click", () => {
    panel.hidden = !panel.hidden;
    if (!panel.hidden) void refresh();
  });

  bus.watchChannel(`user/${userId}/notifications`, () => {
    void refresh();
  });

  void refresh();
}
