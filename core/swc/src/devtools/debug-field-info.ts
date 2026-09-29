import type { SwcArchField } from "../types/workspace.js";
import { SwcRecord } from "../model/record.js";
import { fieldDomain, resolveFieldModifiers } from "../model/modifiers.js";

const SENSITIVE_FIELD_RE =
  /^(password|passwd|pass|token|secret|api_key|apikey|sid|auth|credential|private_key)$/i;

export interface FieldDebugInfo {
  model: string;
  fieldName: string;
  label: string;
  type: string;
  widget: string;
  relation: string;
  staticInvisible: boolean;
  staticReadonly: boolean;
  staticRequired: boolean;
  invisibleExpr: string;
  readonlyExpr: string;
  requiredExpr: string;
  evaluatedInvisible: boolean;
  evaluatedReadonly: boolean;
  evaluatedRequired: boolean;
  aclReadDenied: boolean;
  aclWriteDenied: boolean;
  valueDisplay: string;
  modifierSummary: string;
  domainDisplay: string;
  contextDisplay: string;
  relatedDisplay: string;
}

export function formatDebugDomain(field?: SwcArchField, record?: SwcRecord): string {
  if (!field) return "—";
  const domain = fieldDomain(field, record);
  if (domain === undefined) return "—";
  try {
    const raw = JSON.stringify(domain);
    return raw.length > 200 ? `${raw.slice(0, 197)}…` : raw;
  } catch {
    return "—";
  }
}

export function formatDebugContext(context?: Record<string, unknown>): string {
  if (!context || Object.keys(context).length === 0) return "{}";
  try {
    const raw = JSON.stringify(context);
    return raw.length > 200 ? `${raw.slice(0, 197)}…` : raw;
  } catch {
    return "{}";
  }
}

export function parseDebugFieldKey(key: string): { model: string; fieldName: string } {
  const dot = key.lastIndexOf(".");
  if (dot <= 0) return { model: "", fieldName: key };
  return { model: key.slice(0, dot), fieldName: key.slice(dot + 1) };
}

export function isSensitiveDebugField(fieldName: string): boolean {
  const base = fieldName.split(".").pop() ?? fieldName;
  return SENSITIVE_FIELD_RE.test(base.trim());
}

export function formatDebugFieldValue(fieldName: string, value: unknown): string {
  if (isSensitiveDebugField(fieldName)) return "— (redacted)";
  if (value === null || value === undefined) return "—";
  if (typeof value === "object") {
    try {
      const raw = JSON.stringify(value);
      return raw.length > 120 ? `${raw.slice(0, 117)}…` : raw;
    } catch {
      return String(value);
    }
  }
  const text = String(value);
  return text.length > 120 ? `${text.slice(0, 117)}…` : text;
}

export function modifierSummary(info: Pick<FieldDebugInfo, "evaluatedInvisible" | "evaluatedReadonly" | "evaluatedRequired">): string {
  const parts: string[] = [];
  if (info.evaluatedInvisible) parts.push("invisible");
  if (info.evaluatedReadonly) parts.push("readonly");
  if (info.evaluatedRequired) parts.push("required");
  return parts.length ? parts.join(", ") : "none";
}

export interface BuildFieldDebugInfoInput {
  fieldKey: string;
  archField?: SwcArchField;
  domType?: string;
  domWidget?: string;
  domLabel?: string;
  record?: SwcRecord;
  userId?: number;
  companyId?: number;
  recordData?: Record<string, unknown>;
  aclField?: { read_denied?: boolean; write_denied?: boolean };
  viewContext?: Record<string, unknown>;
}

export function buildFieldDebugInfo(input: BuildFieldDebugInfoInput): FieldDebugInfo {
  const { model, fieldName } = parseDebugFieldKey(input.fieldKey);
  const arch = input.archField;
  const type = (input.domType || arch?.type || "").trim();
  const widget = (input.domWidget || arch?.widget || "").trim();
  const label = (input.domLabel || arch?.string || fieldName).trim();

  let evaluatedInvisible = arch?.invisible ?? false;
  let evaluatedReadonly = arch?.readonly ?? false;
  let evaluatedRequired = arch?.required ?? false;
  if (arch) {
    const mods = resolveFieldModifiers(arch, input.record, {
      userId: input.userId,
      companyId: input.companyId,
    });
    evaluatedInvisible = mods.invisible;
    evaluatedReadonly = mods.readonly;
    evaluatedRequired = mods.required;
  }

  const valueRaw = input.recordData?.[fieldName];
  const info: FieldDebugInfo = {
    model,
    fieldName,
    label,
    type,
    widget,
    relation: (arch?.relation ?? "").trim(),
    staticInvisible: arch?.invisible ?? false,
    staticReadonly: arch?.readonly ?? false,
    staticRequired: arch?.required ?? false,
    invisibleExpr: (arch?.invisible_expr ?? "").trim(),
    readonlyExpr: (arch?.readonly_expr ?? "").trim(),
    requiredExpr: (arch?.required_expr ?? "").trim(),
    evaluatedInvisible,
    evaluatedReadonly,
    evaluatedRequired,
    aclReadDenied: input.aclField?.read_denied === true,
    aclWriteDenied: input.aclField?.write_denied === true,
    valueDisplay: formatDebugFieldValue(fieldName, valueRaw),
    modifierSummary: "",
    domainDisplay: formatDebugDomain(arch, input.record),
    contextDisplay: formatDebugContext(input.viewContext),
    relatedDisplay: (arch?.options?.inverse ?? arch?.options?.relation ?? "").trim(),
  };
  info.modifierSummary = modifierSummary(info);
  return info;
}

/** Bullet lines for field popover (Label, Field, Model, …). */
export function fieldDebugPopoverLines(info: FieldDebugInfo): string[] {
  return [
    `Label: ${info.label}`,
    `Field: ${info.fieldName}`,
    `Model: ${info.model}`,
    `Type: ${info.type || "—"}`,
    `Context: ${info.contextDisplay}`,
    `Domain: ${info.domainDisplay}`,
    `Relation: ${info.relation || "—"}`,
    `Related: ${info.relatedDisplay || "—"}`,
    `Modifiers: ${info.modifierSummary}`,
  ];
}

export function fieldDebugInfoFromElement(
  el: HTMLElement,
  workspace?: { fields?: SwcArchField[]; recordData?: Record<string, unknown>; model?: string; recordId?: number },
  boot?: { userId: number; companyId: number },
  aclTrace?: { fields?: { name: string; read_denied?: boolean; write_denied?: boolean }[] },
): FieldDebugInfo | null {
  const key = el.getAttribute("data-sum-debug-field");
  if (!key) return null;
  const { fieldName } = parseDebugFieldKey(key);
  const archField = workspace?.fields?.find((f) => f.name === fieldName);
  const record =
    workspace?.recordData && workspace.model && workspace.recordId != null
      ? new SwcRecord(workspace.model, workspace.recordId, workspace.recordData)
      : undefined;
  return buildFieldDebugInfo({
    fieldKey: key,
    archField,
    domType: el.getAttribute("data-sum-debug-type") ?? undefined,
    domWidget: el.getAttribute("data-sum-debug-widget") ?? undefined,
    domLabel: el.getAttribute("data-sum-debug-label") ?? undefined,
    record,
    userId: boot?.userId,
    companyId: boot?.companyId,
    recordData: workspace?.recordData,
    aclField: aclTrace?.fields?.find((f) => f.name === fieldName),
    viewContext: boot
      ? { user_id: boot.userId, company_id: boot.companyId }
      : undefined,
  });
}
