import { html, type TemplateResult, type TemplateValue } from "../../template/html.js";
import type { SwcArchField, SwcReportMeta, SwcSearchFilter, SwcWorkspacePayload } from "../../types/workspace.js";
import { inputValueFromEvent } from "../../widgets/field-events.js";
import { buildAllReportEntries, createToolbarIcon, type ReportActionEntry } from "./view-toolbar.js";
import {
  filterOperatorsForField,
  type CollectionQuery,
} from "./collection-query.js";

function renderPopoverShell(modifierClass: string, heading: string, body: TemplateValue): TemplateResult {
  return html`
    <div class=${`sum-popover ${modifierClass}`} @click=${(e: Event) => e.stopPropagation()}>
      <h3 class="sum-popover-heading">${heading}</h3>
      ${body}
    </div>
  `;
}

export function renderPopoverCheckmark(active: boolean): TemplateResult {
  return html`<span class="sum-popover-check" aria-hidden="true">${active ? "✓" : ""}</span>`;
}

export function renderPopoverItem(label: string, active: boolean, onClick: () => void): TemplateResult {
  return html`<li>
    <button
      type="button"
      class=${active ? "sum-popover-item sum-popover-item--active" : "sum-popover-item"}
      @click=${onClick}
    >
      ${renderPopoverCheckmark(active)}
      <span class="sum-popover-item-label">${label}</span>
    </button>
  </li>`;
}

export interface FiltersPopoverState {
  query: CollectionQuery;
  customField: string;
  customOp: string;
  customValue: string;
  domainPresets: SwcSearchFilter[];
  filterFields: SwcArchField[];
}

export interface FiltersPopoverCallbacks {
  onTogglePreset: (name: string) => void;
  onCustomFieldChange: (field: string, defaultOp: string) => void;
  onCustomOpChange: (op: string) => void;
  onCustomValueInput: (value: string) => void;
  onApplyCustom: () => void;
}

export function renderFiltersPopover(
  state: FiltersPopoverState,
  callbacks: FiltersPopoverCallbacks,
): TemplateResult {
  const { query, customField, customOp, customValue, domainPresets, filterFields } = state;
  const field = customField || filterFields[0]?.name || "";
  const operators = filterOperatorsForField(field, filterFields);

  return renderPopoverShell(
    "sum-popover--filters",
    "Filters",
    html`
      <ul class="sum-popover-list">
        ${domainPresets.map((f) =>
          renderPopoverItem(
            f.string || f.name,
            query.presetFilters.includes(f.name),
            () => callbacks.onTogglePreset(f.name),
          ),
        )}
      </ul>
      <div class="sum-popover-custom">
        <strong class="sum-popover-custom-title">Custom filter</strong>
        <select
          class="sum-popover-select"
          @change=${(e: Event) => {
            const next = (e.target as HTMLSelectElement).value;
            callbacks.onCustomFieldChange(next, filterOperatorsForField(next, filterFields)[0] ?? "=");
          }}
        >
          ${filterFields.map((f) =>
            html`<option value=${f.name} selected=${f.name === field ? "selected" : undefined}>${f.string || f.name}</option>`,
          )}
        </select>
        <select class="sum-popover-select" @change=${(e: Event) => callbacks.onCustomOpChange((e.target as HTMLSelectElement).value)}>
          ${operators.map((op) =>
            html`<option value=${op} selected=${op === customOp ? "selected" : undefined}>${op}</option>`,
          )}
        </select>
        <input
          type="text"
          class="sum-popover-input"
          placeholder="Value"
          value=${customValue}
          @input=${(e: Event) => callbacks.onCustomValueInput(inputValueFromEvent(e))}
        />
        <button type="button" class="sum-btn sum-btn--secondary" @click=${() => callbacks.onApplyCustom()}>Apply</button>
      </div>
    `,
  );
}

export interface GroupPopoverState {
  query: CollectionQuery;
  groupPresets: SwcSearchFilter[];
  groupByFields: SwcArchField[];
}

export interface GroupPopoverCallbacks {
  onToggleGroupBy: (field: string) => void;
}

export function renderGroupPopover(state: GroupPopoverState, callbacks: GroupPopoverCallbacks): TemplateResult {
  const { query, groupPresets, groupByFields } = state;

  return renderPopoverShell(
    "sum-popover--group",
    "Group By",
    html`<ul class="sum-popover-list">
      ${groupPresets.map((f) =>
        renderPopoverItem(
          f.string || f.name,
          query.groupBy.includes(f.groupBy!),
          () => callbacks.onToggleGroupBy(f.groupBy!),
        ),
      )}
      ${groupByFields.map((f) =>
        renderPopoverItem(
          f.string || f.name,
          query.groupBy.includes(f.name),
          () => callbacks.onToggleGroupBy(f.name),
        ),
      )}
    </ul>`,
  );
}

export interface FavoriteEntry {
  id: number;
  name: string;
  isShared?: boolean;
  search?: string;
  filter?: string;
  domain?: string;
  groupBy?: string;
}

export interface FavoritesPopoverState {
  favorites: FavoriteEntry[];
  saveName: string;
  saveShared: boolean;
  savingFavorite: boolean;
}

export interface FavoritesPopoverCallbacks {
  onApplyFavorite: (fav: FavoriteEntry) => void;
  onDeleteFavorite: (id: number) => void;
  onSaveNameInput: (name: string) => void;
  onSaveSharedChange: (shared: boolean) => void;
  onSaveFavorite: () => void;
}

export function renderFavoritesPopover(
  state: FavoritesPopoverState,
  callbacks: FavoritesPopoverCallbacks,
): TemplateResult {
  const { favorites, saveName, saveShared, savingFavorite } = state;

  return renderPopoverShell(
    "sum-popover--favorites",
    "Favorites",
    html`
      <ul class="sum-popover-list">
        ${favorites.map(
          (f) => html`<li class="sum-popover-fav">
            <button type="button" class="sum-popover-item" @click=${() => callbacks.onApplyFavorite(f)}>
              <span class="sum-popover-item-label">${f.name}${f.isShared ? " (shared)" : ""}</span>
            </button>
            <button type="button" class="sum-popover-fav-delete" @click=${() => void callbacks.onDeleteFavorite(f.id)} title="Delete">×</button>
          </li>`,
        )}
      </ul>
      <div class="sum-popover-custom">
        <input
          type="text"
          class="sum-popover-input"
          placeholder="Favorite name"
          value=${saveName}
          @input=${(e: Event) => callbacks.onSaveNameInput(inputValueFromEvent(e))}
        />
        <label class="sum-popover-save-shared">
          <input
            type="checkbox"
            ?checked=${saveShared}
            @change=${(e: Event) => callbacks.onSaveSharedChange((e.target as HTMLInputElement).checked)}
          />
          Share with all users
        </label>
        <button
          type="button"
          class="sum-btn sum-btn--secondary"
          disabled=${savingFavorite ? "disabled" : undefined}
          @click=${() => void callbacks.onSaveFavorite()}
        >
          Save current search
        </button>
      </div>
    `,
  );
}

function renderReportsMenu(entries: ReportActionEntry[], popoverExtraClass = ""): TemplateResult {
  const popoverClass = ["sum-popover", "sum-popover--actions", "sum-reports-popover", popoverExtraClass]
    .filter(Boolean)
    .join(" ");

  return html`
    <div class=${popoverClass} @click=${(e: Event) => e.stopPropagation()}>
      <h3 class="sum-popover-heading">Import / export</h3>
      ${entries.length > 0
        ? html`<ul class="sum-popover-menu">
            ${entries.map((entry) => html`<li class="sum-popover-menu-item">${entry.node}</li>`)}
          </ul>`
        : html`<p class="sum-popover-empty">No export actions available.</p>`}
    </div>
  `;
}

export function reportsToolbarVisible(report?: SwcReportMeta): boolean {
  return !!report && (report.download || report.upload);
}

export function renderReportsPopover(
  payload: SwcWorkspacePayload,
  fieldsCsv: string,
  viewType?: string,
  recordId = 0,
  popoverExtraClass = "",
): TemplateResult {
  return renderReportsMenu(
    buildAllReportEntries(payload, fieldsCsv, recordId, viewType),
    popoverExtraClass,
  );
}

export interface ReportsAnchorOptions {
  open: boolean;
  onToggle: () => void;
  payload: SwcWorkspacePayload;
  fieldsCsv: string;
  recordId?: number;
  viewType?: string;
  variant: "collection" | "form";
}

/** Reports trigger + popover (collection bar and form toolbar). */
export function renderReportsAnchor(options: ReportsAnchorOptions): TemplateResult | string {
  const { open, onToggle, payload, fieldsCsv, recordId = 0, viewType, variant } = options;
  if (!reportsToolbarVisible(payload.arch.report)) {
    return "";
  }
  const entries = buildAllReportEntries(payload, fieldsCsv, recordId, viewType);

  const popoverExtraClass = variant === "form" ? "sum-form-reports-popover" : "";
  const popover = open ? renderReportsMenu(entries, popoverExtraClass) : "";

  if (variant === "form") {
    return html`
      <div class="sum-view-toolbar-actions sum-reports-anchor sum-form-reports-anchor" @click=${(e: Event) => e.stopPropagation()}>
        <button
          type="button"
          class="sum-header-btn sum-header-btn--secondary sum-reports-btn sum-form-reports-btn"
          aria-label="Reports"
          title="Import and export records"
          aria-expanded=${open ? "true" : "false"}
          @click=${(e: Event) => { e.stopPropagation(); onToggle(); }}
        >
          Reports
          ${createToolbarIcon("chevron", "sum-reports-chevron sum-form-reports-chevron")}
        </button>
        ${popover}
      </div>
    `;
  }

  const btnClass = open
    ? "sum-control-bar-reports-btn sum-control-bar-reports-btn--active sum-control-bar-actions-btn sum-control-bar-actions-btn--active"
    : "sum-control-bar-reports-btn sum-control-bar-actions-btn";

  return html`
    <div class="sum-control-bar-popover-anchor sum-reports-anchor" @click=${(e: Event) => e.stopPropagation()}>
      <button
        type="button"
        class=${btnClass}
        aria-label="Reports"
        title="Import and export records"
        aria-expanded=${open ? "true" : "false"}
        @click=${(e: Event) => { e.stopPropagation(); onToggle(); }}
      >
        <span class="sum-control-bar-reports-label sum-control-bar-actions-label">Reports</span>
        ${createToolbarIcon("chevron", "sum-control-bar-chip-chevron")}
      </button>
      ${popover}
    </div>
  `;
}
