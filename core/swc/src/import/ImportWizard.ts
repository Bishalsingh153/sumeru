import { SWC_API_BASE } from "../constants/routes.js";

interface ImportBootstrap {
  batchId: number;
  targetModel: string;
  importMode: string;
  headers: string[];
  mapping: Record<string, string>;
  fields: string[];
  fieldLabels: Record<string, string>;
  nextUrl?: string;
}

export function mountImportWizard(host: HTMLElement | null, batchId: string, csrf: string): void {
  if (!host || !batchId) return;
  void load(host, batchId, csrf);
}

async function load(host: HTMLElement, batchId: string, csrf: string): Promise<void> {
  const res = await fetch(`${SWC_API_BASE}/import?batch=${encodeURIComponent(batchId)}`, { credentials: "same-origin" });
  if (!res.ok) {
    host.textContent = "Could not load import batch.";
    return;
  }
  const data = (await res.json()) as ImportBootstrap;
  render(host, data, csrf);
}

function render(host: HTMLElement, data: ImportBootstrap, csrf: string): void {
  host.innerHTML = "";
  const mapping = { ...data.mapping };
  const form = document.createElement("div");
  form.className = "sum-import-wizard";

  const table = document.createElement("table");
  table.className = "sum-import-map-table";
  const thead = document.createElement("thead");
  thead.innerHTML = "<tr><th>File column</th><th>Field</th></tr>";
  table.appendChild(thead);
  const tbody = document.createElement("tbody");
  for (const col of data.headers) {
    const tr = document.createElement("tr");
    const tdCol = document.createElement("td");
    tdCol.textContent = col;
    const tdField = document.createElement("td");
    const select = document.createElement("select");
    select.dataset.column = col;
    appendOption(select, "-", "Ignore");
    for (const f of data.fields) {
      appendOption(select, f, data.fieldLabels[f] || f);
    }
    select.value = mapping[col] ?? "-";
    select.addEventListener("change", () => {
      mapping[col] = select.value;
    });
    tdField.appendChild(select);
    tr.append(tdCol, tdField);
    tbody.appendChild(tr);
  }
  table.appendChild(tbody);
  form.appendChild(table);

  const previewBox = document.createElement("pre");
  previewBox.className = "sum-import-preview";
  form.appendChild(previewBox);

  const actions = document.createElement("div");
  actions.className = "sum-import-actions";

  const previewBtn = button("Preview", () => void runPreview(previewBox, data.batchId, mapping, csrf));
  const dryBtn = button("Dry run", () => void runDryRun(previewBox, data.batchId, mapping, csrf));
  const saveTplBtn = button("Save template", () => void runSaveTemplate(previewBox, data.batchId, mapping, csrf));
  const importBtn = button("Import", () => void runImport(data.batchId, mapping, csrf, data.nextUrl ?? "/web/home"));
  actions.append(previewBtn, dryBtn, saveTplBtn, importBtn);
  form.appendChild(actions);

  host.appendChild(form);
}

function appendOption(select: HTMLSelectElement, value: string, label: string): void {
  const opt = document.createElement("option");
  opt.value = value;
  opt.textContent = label;
  select.appendChild(opt);
}

function button(label: string, onClick: () => void): HTMLButtonElement {
  const b = document.createElement("button");
  b.type = "button";
  b.className = "sum-btn sum-btn--secondary";
  b.textContent = label;
  b.addEventListener("click", onClick);
  return b;
}

async function runPreview(box: HTMLElement, batchId: number, mapping: Record<string, string>, csrf: string): Promise<void> {
  box.textContent = "Previewing…";
  const res = await postJSON("/web/bulk/preview", { batch_id: batchId, column_mapping: mapping }, csrf);
  box.textContent = JSON.stringify(res, null, 2);
}

async function runSaveTemplate(box: HTMLElement, batchId: number, mapping: Record<string, string>, csrf: string): Promise<void> {
  const name = window.prompt("Template name");
  if (!name?.trim()) {
    return;
  }
  box.textContent = "Saving template…";
  const res = await postJSON("/web/bulk/save-template", {
    name: name.trim(),
    batch_id: batchId,
    column_mapping: mapping,
    shared: false,
  }, csrf);
  box.textContent = `Saved template #${String((res as { template_id?: number }).template_id ?? "")}`;
}

async function runDryRun(box: HTMLElement, batchId: number, mapping: Record<string, string>, csrf: string): Promise<void> {
  box.textContent = "Dry run…";
  const res = await postJSON("/web/bulk/dry-run", { batch_id: batchId, column_mapping: mapping }, csrf);
  box.textContent = JSON.stringify(res, null, 2);
}

async function runImport(batchId: number, mapping: Record<string, string>, csrf: string, nextUrl: string): Promise<void> {
  const body = new URLSearchParams({
    id: String(batchId),
    column_mapping: JSON.stringify(mapping),
    csrf_token: csrf,
  });
  const res = await fetch("/web/bulk/confirm", { method: "POST", credentials: "same-origin", body });
  if (res.redirected) {
    window.location.href = res.url;
    return;
  }
  if (res.ok) {
    window.location.href = nextUrl || "/web/home";
  }
}

async function postJSON(url: string, payload: Record<string, unknown>, csrf: string): Promise<unknown> {
  const res = await fetch(url, {
    method: "POST",
    credentials: "same-origin",
    headers: {
      "Content-Type": "application/json",
      Accept: "application/json",
      "X-CSRF-Token": csrf,
    },
    body: JSON.stringify({ ...payload, csrf_token: csrf }),
  });
  if (!res.ok) {
    throw new Error(await res.text());
  }
  return res.json();
}
