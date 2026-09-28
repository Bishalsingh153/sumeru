import type { DialogService } from "../services/dialog.js";
import type { DebugWorkspaceContext } from "./debug.js";

export interface RecordMetadataResponse {
  id?: number;
  create_uid?: unknown;
  write_uid?: unknown;
  create_date?: unknown;
  write_date?: unknown;
  xml_ids?: { xml_id: string }[];
}

function formatMetaValue(raw: unknown): string {
  if (raw === null || raw === undefined || raw === "") return "—";
  if (typeof raw === "object") {
    try {
      return JSON.stringify(raw);
    } catch {
      return String(raw);
    }
  }
  const text = String(raw);
  if (/invalid/i.test(text)) return "—";
  const d = Date.parse(text);
  if (!Number.isNaN(d) && text.includes("-")) {
    try {
      return new Date(d).toLocaleString();
    } catch {
      return text;
    }
  }
  return text;
}

function kvModalBody(pairs: Array<[string, string]>): HTMLElement {
  const wrap = document.createElement("div");
  wrap.className = "sum-debug-metadata";
  for (const [label, value] of pairs) {
    const row = document.createElement("div");
    row.className = "sum-debug-kv";
    const dt = document.createElement("div");
    dt.className = "sum-debug-kv__label";
    dt.textContent = label;
    const dd = document.createElement("div");
    dd.className = "sum-debug-kv__value";
    dd.textContent = value;
    row.append(dt, dd);
    wrap.appendChild(row);
  }
  return wrap;
}

export async function openRecordMetadataModal(
  dialog: DialogService,
  workspace: DebugWorkspaceContext | undefined,
): Promise<void> {
  if (!workspace?.model || !workspace.recordId) {
    await dialog.alert("Metadata", "Open a saved record to view metadata.");
    return;
  }
  let meta: RecordMetadataResponse;
  try {
    const url = `/web/debug/metadata?model=${encodeURIComponent(workspace.model)}&record_id=${workspace.recordId}`;
    const res = await fetch(url, { credentials: "same-origin", headers: { Accept: "application/json" } });
    if (!res.ok) {
      await dialog.alert("Metadata", `Could not load metadata (HTTP ${res.status}).`);
      return;
    }
    meta = (await res.json()) as RecordMetadataResponse;
  } catch {
    await dialog.alert("Metadata", "Could not load metadata.");
    return;
  }

  const xmlLines = (meta.xml_ids ?? []).map((x) => x.xml_id).join("\n") || "—";
  const body = kvModalBody([
    ["ID", String(meta.id ?? workspace.recordId)],
    ["XML ID", xmlLines],
    ["Creation user", formatMetaValue(meta.create_uid)],
    ["Creation date", formatMetaValue(meta.create_date)],
    ["Last modified by", formatMetaValue(meta.write_uid)],
    ["Last modified", formatMetaValue(meta.write_date)],
  ]);
  void dialog.openHost("Metadata", body);
}

export function openRecordDataModal(
  dialog: DialogService,
  workspace: DebugWorkspaceContext | undefined,
  redact: (field: string, value: unknown) => string,
): void {
  if (!workspace?.recordData) {
    void dialog.alert("Data", "No record data loaded for this view.");
    return;
  }
  const lines: Array<[string, string]> = Object.entries(workspace.recordData).map(([k, v]) => [
    k,
    redact(k, v),
  ]);
  void dialog.openHost("Data", kvModalBody(lines));
}
