import { fieldDebugInfoFromElement } from "./debug-field-info.js";
import type { SwcArchField } from "../types/workspace.js";
import type { SwcBootstrap } from "../types/bootstrap.js";

export interface DebugHoverWorkspace {
  model?: string;
  recordId?: number;
  fields?: SwcArchField[];
  recordData?: Record<string, unknown>;
}

const HOVER_CARD_ID = "sum-debug-hover-card";
const THROTTLE_MS = 48;

let hoverEl: HTMLElement | null = null;
let cardEl: HTMLElement | null = null;
let lastMove = 0;
let bootRef: SwcBootstrap | null = null;
let workspaceRef: DebugHoverWorkspace | undefined;
let aclFields: { name: string; read_denied?: boolean; write_denied?: boolean }[] | undefined;
let isDebugActive: () => boolean = () => false;
let hoverListenersBound = false;

export function setDebugHoverContext(
  boot: SwcBootstrap | null,
  workspace?: DebugHoverWorkspace,
  acl?: { fields?: { name: string; read_denied?: boolean; write_denied?: boolean }[] },
): void {
  bootRef = boot;
  workspaceRef = workspace;
  aclFields = acl?.fields;
}

function ensureHoverCard(): HTMLElement {
  if (cardEl?.isConnected) return cardEl;
  cardEl = document.createElement("div");
  cardEl.id = HOVER_CARD_ID;
  cardEl.className = "sum-debug-hover-card";
  cardEl.hidden = true;
  cardEl.setAttribute("role", "status");
  document.body.appendChild(cardEl);
  return cardEl;
}

function hideHoverCard(): void {
  if (cardEl) cardEl.hidden = true;
}

function positionHoverCard(clientX: number, clientY: number): void {
  if (!cardEl || cardEl.hidden) return;
  const pad = 14;
  const rect = cardEl.getBoundingClientRect();
  let left = clientX + pad;
  let top = clientY + pad;
  if (left + rect.width > window.innerWidth - 8) left = clientX - rect.width - pad;
  if (top + rect.height > window.innerHeight - 8) top = clientY - rect.height - pad;
  cardEl.style.left = `${Math.max(8, left)}px`;
  cardEl.style.top = `${Math.max(8, top)}px`;
}

function renderHoverCard(el: HTMLElement, clientX: number, clientY: number): void {
  const card = ensureHoverCard();
  const info = fieldDebugInfoFromElement(
    el,
    workspaceRef,
    bootRef ? { userId: bootRef.user.id, companyId: bootRef.activeCompanyId } : undefined,
    { fields: aclFields },
  );
  if (!info) {
    card.hidden = true;
    return;
  }
  card.replaceChildren();
  const title = document.createElement("div");
  title.className = "sum-debug-hover-card__title";
  title.textContent = info.fieldName;
  const sub = document.createElement("div");
  sub.className = "sum-debug-hover-card__sub";
  sub.textContent = [info.label, info.type, info.widget].filter(Boolean).join(" · ");
  const mods = document.createElement("div");
  mods.className = "sum-debug-hover-card__mods";
  mods.textContent = `Modifiers: ${info.modifierSummary}`;
  card.append(title, sub, mods);
  card.hidden = false;
  positionHoverCard(clientX, clientY);
  requestAnimationFrame(() => positionHoverCard(clientX, clientY));
}

function clearHoverTarget(): void {
  if (hoverEl) hoverEl.classList.remove("sum-debug-field--hover");
  hoverEl = null;
  hideHoverCard();
}

export function setDebugHoverTarget(el: HTMLElement | null, clientX: number, clientY: number): void {
  if (el === hoverEl) {
    if (el) renderHoverCard(el, clientX, clientY);
    return;
  }
  if (hoverEl) hoverEl.classList.remove("sum-debug-field--hover");
  hoverEl = el;
  if (el) {
    el.classList.add("sum-debug-field--hover");
    renderHoverCard(el, clientX, clientY);
  } else {
    hideHoverCard();
  }
}

export function initDebugHoverTrace(active: () => boolean): void {
  if (typeof document === "undefined") return;
  isDebugActive = active;
  ensureHoverCard();
  if (hoverListenersBound) return;
  hoverListenersBound = true;
  document.addEventListener("pointermove", onPointerMove, { passive: true });
}

function onPointerMove(e: PointerEvent): void {
  if (!isDebugActive()) {
    clearHoverTarget();
    return;
  }
  const now = performance.now();
  if (now - lastMove < THROTTLE_MS) return;
  lastMove = now;
  const target = (e.target as Element | null)?.closest("[data-sum-debug-field]") as HTMLElement | null;
  setDebugHoverTarget(target, e.clientX, e.clientY);
}

/** Test hook: reset module state. */
export function resetDebugHoverForTests(): void {
  clearHoverTarget();
  lastMove = 0;
}
