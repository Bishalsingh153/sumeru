import { RECORD_UPDATED, SWC_API_BASE } from "../constants/routes.js";
import type { SwcEnv } from "../runtime/env.js";
import { ACTIVITY_CONTEXT, type ActivityContextPayload } from "./activity-context.js";
import { setActivityLogTabVisible } from "./activity-panel.js";

interface ActivityLogItem {
  meta: string;
  body: string;
}

interface ActivityLogResponse {
  items: ActivityLogItem[];
  hasLog?: boolean;
}

export function initActivityLog(env: SwcEnv): void {
  const host = document.getElementById("sum-activity-pane-log");
  if (!host) return;

  let current: ActivityContextPayload | null = null;

  const clearLogPane = (): void => {
    host.innerHTML = "";
    host.hidden = true;
    setActivityLogTabVisible(false);
  };

  const renderItems = (items: ActivityLogItem[]): void => {
    const rows = items
      .map(
        (item) => `<div class="sum-activity-entry">
        <span class="sum-activity-dot"></span>
        <div><div class="sum-activity-meta">${escapeHtml(item.meta)}</div><div class="sum-activity-msg">${escapeHtml(item.body)}</div></div>
      </div>`,
      )
      .join("");
    host.innerHTML = `<div class="sum-activity-body">${rows}</div>`;
  };

  const loadLog = async (ctx: ActivityContextPayload): Promise<void> => {
    if (!ctx.recordLogEligible || ctx.recordId <= 0 || !ctx.model) {
      clearLogPane();
      return;
    }
    try {
      const base = env.bootstrap.swcApiBase || SWC_API_BASE;
      const data = await env.services.http.getJSON<ActivityLogResponse>(
        `${base}/activity-log?model=${encodeURIComponent(ctx.model)}&id=${ctx.recordId}`,
      );
      const items = data.items ?? [];
      const hasLog = data.hasLog === true || items.length > 0;
      if (!hasLog) {
        clearLogPane();
        return;
      }
      setActivityLogTabVisible(true);
      renderItems(items);
    } catch {
      clearLogPane();
    }
  };

  env.services.bus.subscribe(ACTIVITY_CONTEXT, (payload) => {
    current = payload as ActivityContextPayload;
    void loadLog(current);
  });

  env.services.bus.subscribe(RECORD_UPDATED, (payload) => {
    const msg = payload as { model?: string; id?: number; recordId?: number };
    if (!current?.recordLogEligible) return;
    const rid = msg.id ?? msg.recordId;
    if (msg.model !== current.model || rid !== current.recordId) return;
    void loadLog(current);
  });

  clearLogPane();
}

function escapeHtml(raw: string): string {
  return raw
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;");
}
