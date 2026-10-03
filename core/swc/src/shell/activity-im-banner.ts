import type { SwcBootstrap } from "../types/bootstrap.js";

const BANNER_ID = "sum-im-disabled-banner";

/** Top-of-workspace banner when internal chat is off for the active company (P1). */
export function initInternalMessagesBanner(boot: SwcBootstrap): void {
  if (boot.internalMessagesEnabled !== false) return;
  const content = document.querySelector(".sum-content");
  if (!content || document.getElementById(BANNER_ID)) return;

  const bar = document.createElement("div");
  bar.id = BANNER_ID;
  bar.className = "sum-flash-prem sum-flash-prem--warning";
  bar.setAttribute("role", "status");

  const body = document.createElement("p");
  body.className = "sum-flash-body";
  body.textContent =
    "Internal messages are turned off for this company. Enable Internal chat on the company record to use the Messages panel.";

  const actions = document.createElement("p");
  actions.className = "sum-flash-body";
  if (boot.companyFormHref) {
    const link = document.createElement("a");
    link.href = boot.companyFormHref;
    link.className = "sum-chatter-attachment-link";
    link.textContent = "Open company settings";
    actions.appendChild(link);
  }

  const dismiss = document.createElement("button");
  dismiss.type = "button";
  dismiss.className = "sum-flash-dismiss";
  dismiss.setAttribute("aria-label", "Dismiss");
  dismiss.textContent = "×";
  dismiss.addEventListener("click", () => bar.remove());

  bar.appendChild(dismiss);
  bar.appendChild(body);
  if (boot.companyFormHref) bar.appendChild(actions);
  content.prepend(bar);
}
