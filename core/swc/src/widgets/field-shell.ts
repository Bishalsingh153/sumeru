import { html, type TemplateResult } from "../template/html.js";
import { debugFieldTitle, isDebugMode } from "../devtools/debug.js";
import type { SwcArchField } from "../types/workspace.js";
import type { SwcRecord } from "../model/record.js";

export function fieldInputId(field: SwcArchField): string {
  return `f-${field.name}`;
}

export function fieldLabelId(field: SwcArchField): string {
  return `${fieldInputId(field)}-label`;
}

/** HTML autocomplete token for recognized field names; ERP fields default to off. */
export function fieldAutocomplete(field: SwcArchField): string {
  const explicit = field.options?.autocomplete;
  if (typeof explicit === "string" && explicit.trim()) {
    return explicit.trim();
  }

  const name = field.name.toLowerCase();
  const widget = (field.widget ?? "").toLowerCase();

  if (widget === "email" || name === "email") return "email";
  if (name === "phone" || name === "mobile" || name === "phone_number") return "tel";
  if (name === "name" || name === "display_name") return "organization";
  if (name === "street" || name === "street2") return "street-address";
  if (name === "city") return "address-level2";
  if (name === "zip" || name === "postal_code") return "postal-code";
  if (name === "website" || name === "url") return "url";
  if (name === "firstname" || name === "first_name") return "given-name";
  if (name === "lastname" || name === "last_name") return "family-name";

  return "off";
}

function isFullWidthField(field: SwcArchField): boolean {
  if (field.type === "text" || field.widget === "text") return true;
  if (field.type === "one2many" || field.widget === "one2many") return true;
  if (field.widget === "image") return true;
  return false;
}

export function fieldWidgetClass(field: SwcArchField, extra: string[] = []): string {
  const parts = ["sum-field-widget"];
  if (isFullWidthField(field)) {
    parts.push("sum-field-widget--full");
  }
  if (field.type === "many2one" || field.widget === "many2one") {
    parts.push("sum-field-widget--many2one");
  }
  for (const mod of extra) {
    if (mod) parts.push(mod);
  }
  return parts.join(" ");
}

export function renderDebugFieldInfoBtn(): TemplateResult {
  return html`<button
    type="button"
    class="sum-debug-field-info-btn"
    aria-label="Field debug info"
    tabindex="-1"
  ><svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" aria-hidden="true"><circle cx="12" cy="12" r="10"/><path d="M12 16v-4M12 8h.01"/></svg></button>`;
}

export function fieldLabel(
  field: SwcArchField,
  forId?: string,
  row = false,
  labelId = fieldLabelId(field),
  modelName = "",
): TemplateResult {
  const label = field.string ?? field.name;
  const cls = row ? "sum-field-label sum-field-label--row" : "sum-field-label";
  const title = debugFieldTitle(modelName || "?", field.name, field.type ?? field.widget);
  const infoBtn = isDebugMode() && modelName ? renderDebugFieldInfoBtn() : "";
  const text = html`<span class="sum-field-label__text">${label}</span>`;
  if (forId) {
    return title
      ? html`<label class=${cls} id=${labelId} for=${forId} title=${title}>${text}${infoBtn}</label>`
      : html`<label class=${cls} id=${labelId} for=${forId}>${text}${infoBtn}</label>`;
  }
  return title
    ? html`<span class=${cls} id=${labelId} title=${title}>${text}${infoBtn}</span>`
    : html`<span class=${cls} id=${labelId}>${text}${infoBtn}</span>`;
}

export function fieldControl(
  body: TemplateResult | string,
  compact = false,
  ariaLabelledBy?: string,
): TemplateResult {
  const cls = compact ? "sum-field-control sum-field-control--compact" : "sum-field-control";
  if (ariaLabelledBy) {
    return html`<div class=${cls} aria-labelledby=${ariaLabelledBy}>${body}</div>`;
  }
  return html`<div class=${cls}>${body}</div>`;
}

export function fieldPlaceholder(field: SwcArchField): string {
  return field.placeholder ?? field.string ?? field.name;
}

export function fieldReadonlyValue(val: string, placeholder = ""): TemplateResult {
  const hasValue = val.trim() !== "";
  const text = hasValue ? val : placeholder;
  const cls = hasValue ? "sum-field-value" : "sum-field-value sum-field-value--placeholder";
  return html`<div class=${cls}>${text}</div>`;
}

export function fieldReadonlyInput(
  field: SwcArchField,
  val: string,
  inputType = "text",
): TemplateResult {
  const placeholder = fieldPlaceholder(field);
  return html`<input
    type=${inputType}
    id=${fieldInputId(field)}
    class="sum-field-input"
    name=${field.name}
    value=${val}
    placeholder=${placeholder}
    autocomplete=${fieldAutocomplete(field)}
    readonly
    tabindex="-1"
  />`;
}

/** Merge record model into field shell options (pass on every renderFieldShell call). */
export function shellOptions(record: SwcRecord, extra: FieldShellOptions = {}): FieldShellOptions {
  return { ...extra, record, modelName: extra.modelName ?? record.model };
}

export interface FieldShellOptions {
  showLabel?: boolean;
  modifiers?: string[];
  labelFor?: string | false;
  layout?: "row" | "stack";
  compact?: boolean;
  modelName?: string;
  record?: SwcRecord;
}

function debugShellAttrs(modelName: string, field: SwcArchField) {
  const debugKey = modelName ? `${modelName}.${field.name}` : undefined;
  return {
    debugKey,
    type: field.type ?? "",
    widget: field.widget ?? "",
    label: field.string ?? field.name,
  };
}

export function renderFieldShell(
  field: SwcArchField,
  body: TemplateResult | string,
  options: FieldShellOptions = {},
): TemplateResult {
  const showLabel = options.showLabel !== false;
  const labelId = fieldLabelId(field);
  const modelName = options.modelName ?? options.record?.model ?? "";
  const labelFor = options.labelFor === false ? undefined : options.labelFor;
  const useRow =
    options.layout === "row" || (options.layout !== "stack" && !isFullWidthField(field) && !options.compact);
  const modifiers = [...(options.modifiers ?? [])];
  if (useRow) modifiers.push("sum-field-widget--row");
  if (isDebugMode() && modelName) modifiers.push("sum-field-widget--debug");
  const wrappedBody = fieldControl(body, options.compact === true, labelFor ? undefined : labelId);
  const attrs = debugShellAttrs(modelName, field);
  const floatingDebug =
    !showLabel && isDebugMode() && modelName
      ? html`<div class="sum-field-widget__debug-anchor">${renderDebugFieldInfoBtn()}</div>`
      : "";

  const shellInner = html`
    ${floatingDebug}
    ${showLabel ? fieldLabel(field, labelFor, useRow, labelId, modelName) : ""}
    ${wrappedBody}
  `;

  return html`<div
    class=${fieldWidgetClass(field, modifiers)}
    data-sum-debug-field=${attrs.debugKey}
    data-sum-debug-type=${attrs.type}
    data-sum-debug-widget=${attrs.widget}
    data-sum-debug-label=${attrs.label}
  >
    ${shellInner}
  </div>`;
}

/** Debug wrapper for custom form regions (hero, contact row). */
export function wrapDebugFieldRegion(
  modelName: string,
  field: SwcArchField,
  body: TemplateResult,
): TemplateResult {
  if (!isDebugMode() || !modelName) return body;
  const attrs = debugShellAttrs(modelName, field);
  return html`<div
    class="sum-debug-field-region"
    data-sum-debug-field=${attrs.debugKey}
    data-sum-debug-type=${attrs.type}
    data-sum-debug-widget=${attrs.widget}
    data-sum-debug-label=${attrs.label}
  >
    <div class="sum-debug-field-region__anchor">${renderDebugFieldInfoBtn()}</div>
    ${body}
  </div>`;
}
