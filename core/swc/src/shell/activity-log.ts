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
  hasMore?: boolean;
  offset?: number;
}

export function initActivityLog(env: SwcEnv): void {
  const host = document.getElementById("sum-activity-pane-log");
  if (!host) return;

  let current: ActivityContextPayload | null = null;
  let logOffset = 0;
  let logHasMore = false;
  let cachedItems: ActivityLogItem[] = [];

  const clearLogPane = (): void => {
    host.innerHTML = "";
    host.hidden = true;
    setActivityLogTabVisible(false);
    logOffset = 0;
    logHasMore = false;
    cachedItems = [];
  };

  const renderListHint = (): void => {
    host.innerHTML = `<div class="sum-activity-body"><p class="sum-msg-thread-empty-hint">Open a saved record (form view) to view its change history here.</p></div>`;
    host.hidden = false;
  };

  const renderEmpty = (): void => {
    host.innerHTML = `<div class="sum-activity-body"><p class="sum-msg-thread-empty-hint">No changes logged yet.</p></div>`;
    host.hidden = false;
  };

  const renderItems = (items: ActivityLogItem[], append: boolean): void => {
    if (items.length === 0 && !append) {
      renderEmpty();
      return;
    }
    const all = append ? [...cachedItems, ...items] : items;
    cachedItems = all;
    const rows = all
      .map(
        (item) => `<div class="sum-activity-entry">
        <span class="sum-activity-dot"></span>
        <div><div class="sum-activity-meta">${escapeHtml(item.meta)}</div><div class="sum-activity-msg">${escapeHtml(item.body)}</div></div>
      </div>`,
      )
      .join("");
    const moreBtn = logHasMore
      ? `<button type="button" class="sum-btn sum-btn--secondary sum-activity-log-more">Load more</button>`
      : "";
    host.innerHTML = `<div class="sum-activity-body">${rows}${moreBtn}</div>`;
    host.hidden = false;
    host.querySelector(".sum-activity-log-more")?.addEventListener("click", () => {
      if (current) void loadLog(current, logOffset, true);
    });
  };

  const loadLog = async (ctx: ActivityContextPayload, offset = 0, append = false): Promise<void> => {
    if (!ctx.recordLogEligible || !ctx.model) {
      clearLogPane();
      return;
    }
    setActivityLogTabVisible(true);
    if (ctx.recordId <= 0) {
      renderListHint();
      return;
    }
    try {
      const base = env.bootstrap.swcApiBase || SWC_API_BASE;
      const data = await env.services.http.getJSON<ActivityLogResponse>(
        `${base}/activity-log?model=${encodeURIComponent(ctx.model)}&id=${ctx.recordId}&offset=${offset}`,
      );
      logOffset = (data.offset ?? offset) + (data.items?.length ?? 0);
      logHasMore = data.hasMore === true;
      const items = data.items ?? [];
      if (data.hasLog === true || items.length > 0 || append) {
        renderItems(items, append);
      } else {
        renderEmpty();
      }
    } catch {
      renderEmpty();
    }
  };

  env.services.bus.subscribe(ACTIVITY_CONTEXT, (payload) => {
    current = payload as ActivityContextPayload;
    logOffset = 0;
    logHasMore = false;
    cachedItems = [];
    void loadLog(current);
  });

  env.services.bus.subscribe(RECORD_UPDATED, (payload) => {
    const msg = payload as { model?: string; id?: number; recordId?: number };
    if (!current?.recordLogEligible) return;
    const rid = msg.id ?? msg.recordId;
    if (msg.model !== current.model || rid !== current.recordId) return;
    logOffset = 0;
    logHasMore = false;
    cachedItems = [];
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
