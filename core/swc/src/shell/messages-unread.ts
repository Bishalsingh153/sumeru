import type { HttpService } from "../services/http.js";
import { SWC_API_BASE } from "../constants/routes.js";
import { setMessagesUnreadBadge } from "./activity-panel.js";

interface DirectStatusPayload {
  enabled?: boolean;
  unread?: number;
}

export function initMessagesUnreadPoll(http: HttpService, apiBase?: string): void {
  const base = apiBase || SWC_API_BASE;

  const refresh = async (): Promise<void> => {
    try {
      const data = await http.getJSON<DirectStatusPayload>(`${base}/direct/status`);
      if (data.enabled === false) {
        setMessagesUnreadBadge(0);
        return;
      }
      setMessagesUnreadBadge(data.unread ?? 0);
    } catch {
      setMessagesUnreadBadge(0);
    }
  };

  void refresh();
  window.setInterval(() => void refresh(), 60_000);
}
