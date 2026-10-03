import type { SwcEnv } from "../runtime/env.js";
import { InternalMessagesPanel } from "./internal-messages-panel.js";

/** Messages tab: internal user search and P2P chat (not record chatter or Log). */
export function initActivityMessages(env: SwcEnv): void {
  const host = document.getElementById("sum-activity-pane-messages");
  if (!host) return;

  host.innerHTML = "";
  const mountEl = document.createElement("div");
  mountEl.className = "sum-activity-messages-mount";
  host.appendChild(mountEl);

  const panel = new InternalMessagesPanel({}, env);
  panel.callSetup();
  mountEl.replaceChildren(panel.render());
}
