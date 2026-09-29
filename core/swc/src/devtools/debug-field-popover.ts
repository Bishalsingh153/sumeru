import { fieldDebugInfoFromElement, fieldDebugPopoverLines } from "./debug-field-info.js";
import {
  getDebugMenuContext,
  isFieldInspectorEnabled,
  openDebugDrawer,
  selectDebugFieldFromElement,
} from "./debug.js";

const POPOVER_ID = "sum-debug-field-popover";

function ensurePopover(): HTMLElement {
  let el = document.getElementById(POPOVER_ID);
  if (el) return el;
  el = document.createElement("div");
  el.id = POPOVER_ID;
  el.className = "sum-debug-field-popover";
  el.hidden = true;
  el.setAttribute("role", "tooltip");
  document.body.appendChild(el);
  return el;
}

function fieldWidgetFromIcon(btn: HTMLElement): HTMLElement | null {
  return btn.closest("[data-sum-debug-field]") as HTMLElement | null;
}

function positionPopover(anchor: HTMLElement): void {
  const pop = ensurePopover();
  const rect = anchor.getBoundingClientRect();
  pop.style.left = `${Math.min(rect.left, window.innerWidth - 290)}px`;
  pop.style.top = `${rect.bottom + 8}px`;
}

export function showFieldDebugPopover(anchor: HTMLElement): void {
  const widget = fieldWidgetFromIcon(anchor) ?? anchor.closest("[data-sum-debug-field]");
  if (!widget) return;
  const ctx = getDebugMenuContext();
  const info = fieldDebugInfoFromElement(
    widget as HTMLElement,
    ctx.workspace,
    ctx.boot ? { userId: ctx.boot.user.id, companyId: ctx.boot.activeCompanyId } : undefined,
    ctx.aclTrace ?? undefined,
  );
  if (!info) return;
  const pop = ensurePopover();
  pop.replaceChildren();
  const ul = document.createElement("ul");
  ul.className = "sum-debug-field-popover__list";
  for (const line of fieldDebugPopoverLines(info)) {
    const li = document.createElement("li");
    li.textContent = line;
    ul.appendChild(li);
  }
  pop.appendChild(ul);
  pop.hidden = false;
  widget.classList.add("sum-debug-field--popover");
  positionPopover(anchor);
}

export function hideFieldDebugPopover(): void {
  const pop = document.getElementById(POPOVER_ID);
  if (pop) pop.hidden = true;
  document.querySelectorAll(".sum-debug-field--popover").forEach((el) => {
    el.classList.remove("sum-debug-field--popover");
  });
}

export function initDebugFieldPopovers(): void {
  if (typeof document === "undefined") return;
  ensurePopover();

  document.addEventListener(
    "mouseover",
    (e) => {
      const btn = (e.target as Element).closest(".sum-debug-field-info-btn") as HTMLElement | null;
      if (btn) showFieldDebugPopover(btn);
    },
    true,
  );

  document.addEventListener(
    "mouseout",
    (e) => {
      const related = (e as MouseEvent).relatedTarget as Element | null;
      if (related?.closest?.(".sum-debug-field-info-btn") || related?.closest?.(`#${POPOVER_ID}`)) return;
      if ((e.target as Element).closest(".sum-debug-field-info-btn")) hideFieldDebugPopover();
    },
    true,
  );

  document.addEventListener(
    "click",
    (e) => {
      const btn = (e.target as Element).closest(".sum-debug-field-info-btn") as HTMLElement | null;
      if (!btn) return;
      e.preventDefault();
      e.stopPropagation();
      const widget = fieldWidgetFromIcon(btn);
      if (!widget) return;
      selectDebugFieldFromElement(widget);
      openDebugDrawer();
      if (isFieldInspectorEnabled()) {
        e.stopImmediatePropagation();
      }
    },
    true,
  );
}

export function selectDebugFieldByName(model: string, fieldName: string): void {
  const key = `${model}.${fieldName}`;
  const el = document.querySelector<HTMLElement>(`[data-sum-debug-field="${key}"]`);
  if (el) selectDebugFieldFromElement(el);
}
