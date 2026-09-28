/** Top-bar developer debug menu (system admins only). */

import type { SwcBootstrap } from "../types/bootstrap.js";
import type { DialogService } from "../services/dialog.js";
import type { NotificationService } from "../services/notification.js";
import { buildDebugMenuSections } from "./debug-menu-sections.js";
import {
  applyDebugBodyClasses,
  getDebugMenuContext,
  isDebugMode,
  isFieldInspectorEnabled,
  logViewArch,
  navigateDebugMode,
  onDebugSelectionChange,
  openDebugDrawer,
  refreshDebugDrawer,
  registerDebugMenuItem,
  selectDebugFieldFromElement,
  toggleFieldInspector,
} from "./debug.js";
import { initDebugFieldPopovers } from "./debug-field-popover.js";
import { formatDebugFieldValue } from "./debug-field-info.js";
import { openRecordDataModal, openRecordMetadataModal } from "./debug-metadata.js";
import { mountDevtoolsPanel } from "./panel.js";

const BUG_ICON = `<svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><path d="M8 2v4M16 2v4M12 2v2M5 10a7 7 0 0 1 14 0c0 4-2 6-2 9H7c0-3-2-5-2-9z"/><path d="M9 17v2M15 17v2M12 17v3"/></svg>`;

export interface DebugManagerServices {
  dialog: DialogService;
  notification: NotificationService;
}

let servicesRef: DebugManagerServices | null = null;

function closeOpenDropdowns(except?: HTMLDetailsElement): void {
  document.querySelectorAll<HTMLDetailsElement>(".sum-dropdown.sum-dropdown--debug[open]").forEach((el) => {
    if (el !== except) el.open = false;
  });
}

function syncDebugMenuButtonState(btn: HTMLElement): void {
  btn.classList.toggle("sum-debug-menu-btn--active", isDebugMode());
  btn.setAttribute("aria-pressed", isDebugMode() ? "true" : "false");
}

function focusChatterTab(tab: "messages" | "attachments"): void {
  const label = tab === "messages" ? "Messages" : "Attachments";
  const btn = Array.from(document.querySelectorAll<HTMLButtonElement>(".sum-chatter-tab")).find((b) =>
    b.textContent?.startsWith(label),
  );
  if (btn) btn.click();
  else servicesRef?.notification.warning("Chatter", "Not available on this view.");
}

function renderMegaMenu(panel: HTMLDivElement, details: HTMLDetailsElement): void {
  panel.replaceChildren();
  panel.classList.add("sum-debug-mega-menu");

  const addHeading = (text: string) => {
    const h = document.createElement("div");
    h.className = "sum-debug-mega-menu__heading";
    h.textContent = text;
    panel.appendChild(h);
  };

  const addItem = (label: string, run: () => void, opts?: { checked?: boolean; danger?: boolean; readout?: boolean }) => {
    if (opts?.readout) {
      const row = document.createElement("div");
      row.className = "sum-debug-mega-menu__readout";
      row.textContent = label;
      panel.appendChild(row);
      return;
    }
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "sum-debug-mega-menu__item";
    btn.setAttribute("role", "menuitem");
    if (opts?.danger) btn.classList.add("sum-debug-mega-menu__item--danger");
    if (opts?.checked !== undefined) {
      const mark = document.createElement("span");
      mark.className = "sum-dropdown-check";
      mark.textContent = opts.checked ? "✓" : "";
      btn.appendChild(mark);
    }
    btn.appendChild(document.createTextNode(label));
    btn.addEventListener("click", () => {
      run();
      details.open = false;
    });
    panel.appendChild(btn);
  };

  const addSep = () => {
    const hr = document.createElement("div");
    hr.className = "sum-dropdown-sep";
    panel.appendChild(hr);
  };

  if (!isDebugMode()) {
    addItem("Enable developer mode", () => navigateDebugMode("1"));
    addItem("Enable developer mode (assets)", () => navigateDebugMode("assets"));
    return;
  }

  addItem("Disable developer mode", () => navigateDebugMode("off"), { danger: true });
  addSep();

  const { workspace } = getDebugMenuContext();
  if (!servicesRef) return;
  const sections = buildDebugMenuSections(true, workspace);
  const dialog = servicesRef.dialog;

  for (const sec of sections) {
    addHeading(sec.title);
    for (const item of sec.items) {
      if (item.startsWith("Model:") || item.startsWith("View:")) {
        addItem(item, () => {}, { readout: true });
        continue;
      }
      switch (item) {
        case "Metadata":
          addItem(item, () => void openRecordMetadataModal(dialog, workspace));
          break;
        case "Data":
          addItem(item, () => openRecordDataModal(dialog, workspace, formatDebugFieldValue));
          break;
        case "Messages":
          addItem(item, () => focusChatterTab("messages"));
          break;
        case "Attachments":
          addItem(item, () => focusChatterTab("attachments"));
          break;
        case "Set default values":
          addItem(item, () => {
            console.debug("[SWC] Set defaults via server/form defaults");
            servicesRef?.notification.show({ kind: "info", title: "Defaults", body: "Use form or server default values." });
          });
          break;
        case "Fields":
          addItem(item, () => {
            openDebugDrawer("fields");
            void refreshDebugDrawer();
          });
          break;
        case "Filters": {
          addItem(item, () => {
            const fav = workspace?.favorites?.map((f) => f.name).join(", ") || workspace?.listSearch || "—";
            void dialog.alert("Filters", fav);
          });
          break;
        }
        case "View arch":
          addItem(item, () => {
            if (workspace?.arch) logViewArch(workspace.arch);
            const pre = document.createElement("pre");
            pre.className = "sum-debug-trace";
            pre.textContent = workspace?.arch ? JSON.stringify(workspace.arch, null, 2) : "No arch loaded";
            void dialog.openHost("View arch", pre);
          });
          break;
        case "Access rights":
          addItem(item, () => {
            openDebugDrawer("acl");
            void refreshDebugDrawer();
          });
          break;
        case "Regenerate assets":
          addItem(item, () => navigateDebugMode("assets"));
          break;
        case "Field inspector":
          addItem(item, () => {
            toggleFieldInspector();
            applyDebugBodyClasses();
            void refreshDebugDrawer();
          }, { checked: isFieldInspectorEnabled() });
          break;
        case "Open debug drawer":
          addItem(item, () => openDebugDrawer());
          break;
        case "Open SWC Vision":
          addItem(item, () => mountDevtoolsPanel());
          break;
        case "Open metrics":
          addItem(item, () => window.open("/web/metrics", "_blank", "noopener,noreferrer"));
          break;
        case "Reload page":
          addItem(item, () => window.location.reload());
          break;
        default:
          break;
      }
    }
  }
}

export function initDebugManager(boot: SwcBootstrap, services?: DebugManagerServices): void {
  if (!boot.features?.debugMenu) return;
  if (services) servicesRef = services;

  registerDebugMenuItem({
    id: "open-vision",
    label: "Open SWC Vision",
    section: "tools",
    run: () => mountDevtoolsPanel(),
  });
  registerDebugMenuItem({
    id: "log-arch",
    label: "Log view arch to console",
    section: "tools",
    run: () => console.debug("[SWC] Use logViewArch from view load hooks"),
  });

  const host = document.getElementById("sum-debug-menu-host");
  if (!host) return;
  host.hidden = false;
  host.removeAttribute("aria-hidden");

  const details = document.createElement("details");
  details.className = "sum-dropdown sum-dropdown--debug";

  const summary = document.createElement("summary");
  summary.className = "sum-dropdown-trigger sum-icon-btn sum-debug-menu-btn";
  summary.setAttribute("aria-label", "Developer tools");
  summary.setAttribute("title", "Developer tools");
  summary.innerHTML = BUG_ICON;

  const panel = document.createElement("div");
  panel.className = "sum-dropdown-panel sum-dropdown-panel--align-end sum-debug-mega-menu-panel";
  panel.setAttribute("role", "menu");

  details.append(summary, panel);
  host.appendChild(details);

  details.addEventListener("toggle", () => {
    if (details.open) {
      closeOpenDropdowns(details);
      renderMegaMenu(panel, details);
    }
  });

  document.addEventListener("click", (e) => {
    if (!details.open) return;
    if (details.contains(e.target as Node)) return;
    details.open = false;
  });

  syncDebugMenuButtonState(summary);

  document.addEventListener(
    "click",
    (e) => {
      if (!isDebugMode() || !isFieldInspectorEnabled()) return;
      if ((e.target as Element).closest(".sum-debug-field-info-btn")) return;
      const el = (e.target as Element).closest("[data-sum-debug-field]");
      if (!el) return;
      e.preventDefault();
      e.stopPropagation();
      selectDebugFieldFromElement(el as HTMLElement);
    },
    true,
  );

  onDebugSelectionChange(() => {
    void refreshDebugDrawer();
  });

  initDebugFieldPopovers();
  applyDebugBodyClasses();
}
