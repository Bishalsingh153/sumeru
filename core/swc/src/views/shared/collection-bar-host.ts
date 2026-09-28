import { SwcComponent } from "../../runtime/component.js";
import { html, type TemplateResult, type TemplateValue } from "../../template/html.js";
import type { SwcArchField, SwcWorkspacePayload } from "../../types/workspace.js";
import { SAVED_SEARCHES_ROUTE } from "../../constants/routes.js";
import { inputValueFromEvent } from "../../widgets/field-events.js";
import { renderNewButton, buildAllReportEntries, exportFieldNamesCsv, createToolbarIcon } from "./view-toolbar.js";
import { visibleArchFields } from "./arch-fields.js";
import {
  activeFilterTags,
  appendDomainTriple,
  coerceFilterValue,
  filterCount,
  groupByCount,
  navigateCollectionQuery,
  presetDomainFilters,
  removeFilterTag,
  syncCollectionQuery,
  toggleGroupByField,
  togglePresetFilter,
  type CollectionQuery,
  type DomainTriple,
  type FilterTag,
} from "./collection-query.js";
import {
  renderReportsPopover,
  renderFiltersPopover,
  renderFavoritesPopover,
  renderGroupPopover,
} from "./collection-bar-panels.js";

export interface CollectionBarHostProps {
  payload: SwcWorkspacePayload;
  viewType: string;
  extraPrimary?: TemplateValue;
}

type PanelKind = "filters" | "group" | "favorites" | "reports";

const SEGMENT_TOOLTIPS: Record<Exclude<PanelKind, "reports">, string> = {
  filters: "Filter records using presets or custom field rules",
  group: "Group records by field values",
  favorites: "Save and reopen favorite searches",
};

export class CollectionBarHost extends SwcComponent<CollectionBarHostProps> {
  private query!: CollectionQuery;
  private panelOpen: PanelKind | null = null;
  private customField = "";
  private customOp = "=";
  private customValue = "";
  private saveName = "";
  private saveShared = false;
  private savingFavorite = false;

  override setup(): void {
    this.syncQuery();
  }

  override onPropsChanged(props: CollectionBarHostProps): void {
    this.syncQuery();
    void props;
  }

  private syncQuery(): void {
    this.query = syncCollectionQuery(this.props.payload);
  }

  private searchMeta() {
    return this.props.payload.arch.search;
  }

  private applyQuery(patch: Partial<CollectionQuery>): void {
    navigateCollectionQuery(this.env, this.props.payload, this.props.viewType, {
      search: patch.search ?? this.query.search,
      presetFilters: patch.presetFilters ?? this.query.presetFilters,
      customDomain: patch.customDomain ?? this.query.customDomain,
      groupBy: patch.groupBy ?? this.query.groupBy,
      listOffset: 0,
    });
  }

  private applySearch(): void {
    this.applyQuery({ search: this.query.search });
  }

  private togglePanel(kind: PanelKind): void {
    this.panelOpen = this.panelOpen === kind ? null : kind;
    this.rerender();
  }

  private closePanel(): void {
    if (!this.panelOpen) return;
    this.panelOpen = null;
    this.rerender();
  }

  private togglePreset(name: string): void {
    this.applyQuery({ presetFilters: togglePresetFilter(this.query.presetFilters, name) });
  }

  private toggleGroupBy(field: string): void {
    this.applyQuery({ groupBy: toggleGroupByField(this.query.groupBy, field) });
  }

  private applyCustomFilter(): void {
    const field = this.customField.trim();
    if (!field) return;
    const triple: DomainTriple = [field, this.customOp, coerceFilterValue(field, this.customValue, this.filterFields())];
    this.applyQuery({ customDomain: appendDomainTriple(this.query.customDomain, triple) });
    this.customValue = "";
    this.panelOpen = null;
    this.rerender();
  }

  private filterFields(): SwcArchField[] {
    return this.searchMeta()?.filterFields ?? [];
  }

  private groupByFields(): SwcArchField[] {
    const meta = this.searchMeta();
    const fromModel = meta?.groupByFields ?? [];
    const extras = (meta?.filters ?? [])
      .filter((f) => f.groupBy && !fromModel.some((m) => m.name === f.groupBy))
      .map((f) => ({ name: f.groupBy!, string: f.string, type: "char" }));
    return [...fromModel, ...extras];
  }

  private removeTag(tag: FilterTag): void {
    const next = removeFilterTag(this.query, tag);
    this.applyQuery(next);
  }

  private applyFavorite(fav: { search?: string; filter?: string; domain?: string; groupBy?: string }): void {
    this.applyQuery({
      search: fav.search ?? "",
      presetFilters: (fav.filter ?? "").split(",").map((s) => s.trim()).filter(Boolean),
      customDomain: fav.domain ?? "",
      groupBy: (fav.groupBy ?? "").split(",").map((s) => s.trim()).filter(Boolean),
    });
    this.panelOpen = null;
    this.rerender();
  }

  private async saveFavorite(): Promise<void> {
    const name = this.saveName.trim();
    if (!name || this.savingFavorite) return;
    this.savingFavorite = true;
    this.rerender();
    try {
      const payload = this.props.payload;
      await this.env.services.http.postJSON(SAVED_SEARCHES_ROUTE, {
        actionId: payload.actionId,
        model: payload.model,
        name,
        search: this.query.search,
        filter: this.query.presetFilters.join(","),
        domain: this.query.customDomain,
        groupBy: this.query.groupBy.join(","),
        isShared: this.saveShared,
      });
      this.saveName = "";
      this.env.services.notification.success("Saved", "Search saved to favorites.");
      this.panelOpen = null;
    } catch (err) {
      this.env.services.notification.error("Save failed", String(err));
    } finally {
      this.savingFavorite = false;
      this.rerender();
    }
  }

  private async deleteFavorite(id: number): Promise<void> {
    try {
      await this.env.services.http.delete(`${SAVED_SEARCHES_ROUTE}?id=${id}`);
      this.env.services.notification.success("Removed", "Favorite deleted.");
      this.rerender();
    } catch (err) {
      this.env.services.notification.error("Delete failed", String(err));
    }
  }

  private renderBadge(count: number): TemplateResult | string {
    if (count <= 0) return "";
    return html`<span class="sum-control-bar-badge">${count}</span>`;
  }

  private clearSearch(): void {
    this.query.search = "";
    const input = this.rootElement?.querySelector(".sum-control-bar-search") as HTMLInputElement | null;
    if (input) input.value = "";
    if (this.rootElement?.isConnected) this.patch();
    this.applyQuery({ search: "" });
  }

  private renderSearchBar(tags: FilterTag[]): TemplateResult {
    const hasQuery = this.query.search.trim().length > 0;
    return html`
      <div class="sum-control-bar-search-wrap">
        <div class="sum-control-bar-search-body">
          ${tags.length > 0
            ? html`<div class="sum-control-bar-search-tags">${tags.map((tag) => this.renderTag(tag))}</div>`
            : ""}
          <input
            type="search"
            class="sum-control-bar-search"
            placeholder=${tags.length > 0 ? "Search within results…" : "Search records…"}
            aria-label="Search records"
            value=${this.query.search}
            @keydown=${(event: Event) => (event as KeyboardEvent).key === "Enter" && this.applySearch()}
            @input=${(event: Event) => { this.query.search = inputValueFromEvent(event); this.rerender(); }}
            @click=${(e: Event) => e.stopPropagation()}
          />
        </div>
        ${hasQuery
          ? html`<button
              type="button"
              class="sum-control-bar-search-clear"
              aria-label="Clear search"
              title="Clear search"
              @click=${(e: Event) => { e.stopPropagation(); this.clearSearch(); }}
            >
              ${createToolbarIcon("close")}
            </button>`
          : ""}
        <button
          type="button"
          class="sum-control-bar-search-submit"
          aria-label="Search"
          title="Search records"
          @click=${(e: Event) => { e.stopPropagation(); this.applySearch(); }}
        >
          ${createToolbarIcon("search", "sum-control-bar-search-icon-svg")}
        </button>
      </div>
    `;
  }

  private renderTag(tag: FilterTag): TemplateResult {
    return html`<span class="sum-filter-tag sum-filter-tag--${tag.kind}">
      ${tag.label}
      <button type="button" class="sum-filter-tag-remove" aria-label="Remove" @click=${() => this.removeTag(tag)}>×</button>
    </span>`;
  }

  private renderOpenPanel(): TemplateResult | string {
    if (!this.panelOpen) return "";
    const meta = this.searchMeta();
    const filterFields = this.filterFields();
    if (!this.customField && filterFields[0]) this.customField = filterFields[0].name;

    if (this.panelOpen === "filters") {
      return renderFiltersPopover(
        {
          query: this.query,
          customField: this.customField,
          customOp: this.customOp,
          customValue: this.customValue,
          domainPresets: presetDomainFilters(meta?.filters),
          filterFields,
        },
        {
          onTogglePreset: (name) => this.togglePreset(name),
          onCustomFieldChange: (field, defaultOp) => {
            this.customField = field;
            this.customOp = defaultOp;
            this.rerender();
          },
          onCustomOpChange: (op) => { this.customOp = op; },
          onCustomValueInput: (value) => { this.customValue = value; },
          onApplyCustom: () => this.applyCustomFilter(),
        },
      );
    }
    if (this.panelOpen === "group") {
      return renderGroupPopover(
        {
          query: this.query,
          groupPresets: (meta?.filters ?? []).filter((f) => f.groupBy),
          groupByFields: this.groupByFields(),
        },
        { onToggleGroupBy: (field) => this.toggleGroupBy(field) },
      );
    }
    if (this.panelOpen === "favorites") {
      return renderFavoritesPopover(
        {
          favorites: this.props.payload.favorites ?? [],
          saveName: this.saveName,
          saveShared: this.saveShared,
          savingFavorite: this.savingFavorite,
        },
        {
          onApplyFavorite: (fav) => this.applyFavorite(fav),
          onDeleteFavorite: (id) => void this.deleteFavorite(id),
          onSaveNameInput: (name) => { this.saveName = name; },
          onSaveSharedChange: (shared) => { this.saveShared = shared; },
          onSaveFavorite: () => void this.saveFavorite(),
        },
      );
    }
    const fields = exportFieldNamesCsv(visibleArchFields(this.props.payload.arch.fields ?? []));
    return renderReportsPopover(this.props.payload, fields, this.props.viewType);
  }

  private renderFiltersGroupDrop(): TemplateResult | string {
    if (this.panelOpen !== "filters" && this.panelOpen !== "group") {
      return "";
    }
    const dropClass =
      this.panelOpen === "group"
        ? "sum-control-bar-segment-drop sum-control-bar-segment-drop--group"
        : "sum-control-bar-segment-drop sum-control-bar-segment-drop--filters";
    return html`<div class=${dropClass} @click=${(e: Event) => e.stopPropagation()}>
      ${this.renderOpenPanel()}
    </div>`;
  }

  private renderSegmentButton(
    kind: Exclude<PanelKind, "reports">,
    label: string,
    iconName: "filter" | "group" | "favorite",
    badgeCount: number,
  ): TemplateResult {
    const open = this.panelOpen === kind;
    const btnClass = open
      ? "sum-control-bar-segment-btn sum-control-bar-segment-btn--active"
      : "sum-control-bar-segment-btn";
    const badged = badgeCount > 0 ? " sum-control-bar-segment-btn--badged" : "";
    const tooltip = SEGMENT_TOOLTIPS[kind];
    const embedPopover = kind === "favorites";
    return html`
      <div class="sum-control-bar-popover-anchor" data-sum-segment-panel=${kind}>
        <button
          type="button"
          class=${btnClass + badged}
          aria-label=${label}
          title=${tooltip}
          aria-expanded=${open ? "true" : "false"}
          @click=${(e: Event) => { e.stopPropagation(); this.togglePanel(kind); }}
        >
          ${createToolbarIcon(iconName, "sum-control-bar-segment-icon")}
          <span class="sum-control-bar-segment-label">${label}</span>
          ${this.renderBadge(badgeCount)}
        </button>
        ${open && embedPopover ? this.renderOpenPanel() : ""}
      </div>
    `;
  }

  private renderReportsTrigger(entriesCount: number): TemplateResult | string {
    if (entriesCount <= 0) return "";
    const open = this.panelOpen === "reports";
    const btnClass = open
      ? "sum-control-bar-reports-btn sum-control-bar-reports-btn--active sum-control-bar-actions-btn sum-control-bar-actions-btn--active"
      : "sum-control-bar-reports-btn sum-control-bar-actions-btn";
    const chevron = createToolbarIcon("chevron", "sum-control-bar-chip-chevron");
    return html`
      <div class="sum-control-bar-popover-anchor sum-reports-anchor">
        <button
          type="button"
          class=${btnClass}
          aria-label="Reports"
          title="Import and export records"
          aria-expanded=${open ? "true" : "false"}
          @click=${(e: Event) => { e.stopPropagation(); this.togglePanel("reports"); }}
        >
          <span class="sum-control-bar-reports-label sum-control-bar-actions-label">Reports</span>
          ${chevron}
        </button>
        ${open ? this.renderOpenPanel() : ""}
      </div>
    `;
  }

  override template(): TemplateResult {
    const payload = this.props.payload;
    const fields = exportFieldNamesCsv(visibleArchFields(payload.arch.fields ?? []));
    const actionEntries = buildAllReportEntries(payload, fields, 0, this.props.viewType);
    const tags = activeFilterTags(this.query, this.searchMeta());
    const fCount = filterCount(this.query);
    const gCount = groupByCount(this.query);

    return html`
      <div class="sum-control-bar-wrap" @click=${() => this.closePanel()}>
        <aside class="sum-search-panel-facets" aria-label="Search facets">
          ${presetDomainFilters(this.searchMeta()?.filters).map(
            (preset) => {
              const active = this.query.presetFilters.includes(preset.name);
              return html`<button
              type="button"
              class=${active ? "sum-search-facet sum-search-facet--active" : "sum-search-facet"}
              @click=${(e: Event) => {
                e.stopPropagation();
                this.togglePreset(preset.name);
              }}
            >
              ${preset.string || preset.name}
            </button>`;
            },
          )}
        </aside>
        <div class="sum-control-bar sum-view-toolbar">
          <div class="sum-control-bar-row">
            <div class="sum-control-bar-leading">
              ${renderNewButton(payload)}
            </div>
            <div class="sum-control-bar-search-area">
              ${this.renderSearchBar(tags)}
            </div>
            <div class="sum-control-bar-tools">
              <div class="sum-control-bar-segment-wrap">
                <div class="sum-control-bar-segment" role="group" aria-label="Search options">
                  ${this.renderSegmentButton("filters", "Filters", "filter", fCount)}
                  ${this.renderSegmentButton("group", "Group By", "group", gCount)}
                  ${this.renderSegmentButton("favorites", "Favorites", "favorite", 0)}
                </div>
                ${this.renderFiltersGroupDrop()}
              </div>
              ${this.renderReportsTrigger(actionEntries.length)}
              ${this.props.extraPrimary ?? ""}
            </div>
          </div>
        </div>
      </div>
    `;
  }
}

export function mountCollectionBar(
  payload: SwcWorkspacePayload,
  viewType: string,
  env: CollectionBarHost["env"],
  extraPrimary?: TemplateValue,
): CollectionBarHost {
  const host = new CollectionBarHost({ payload, viewType, extraPrimary }, env);
  host.callSetup();
  return host;
}
