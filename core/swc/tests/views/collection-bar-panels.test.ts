import { describe, expect, it, vi } from "vitest";
import {
  renderFavoritesPopover,
  renderFiltersPopover,
  renderGroupPopover,
  renderPopoverItem,
  renderReportsAnchor,
  renderReportsPopover,
  reportsToolbarVisible,
} from "../../src/views/shared/collection-bar-panels.js";
import type { SwcWorkspacePayload } from "../../src/types/workspace.js";

function basePayload(): SwcWorkspacePayload {
  return {
    actionId: 1,
    menuId: "2",
    viewType: "list",
    model: "demo.model",
    recordId: 0,
    formEdit: false,
    csrfToken: "tok",
    arch: {
      type: "list",
      model: "demo.model",
      fields: [{ name: "name" }],
      report: { download: true, formats: "csv" },
    },
    viewTabs: [],
    breadcrumbs: [],
    favorites: [],
  };
}

describe("collection-bar-panels", () => {
  it("renders filter and group popover markup", () => {
    const filters = renderFiltersPopover(
      {
        query: { search: "", presetFilters: ["draft"], customDomain: "", groupBy: [] },
        customField: "name",
        customOp: "=",
        customValue: "",
        domainPresets: [{ name: "draft", string: "Draft", domain: "[]" }],
        filterFields: [{ name: "name", string: "Name", type: "char" }],
      },
      {
        onTogglePreset: vi.fn(),
        onCustomFieldChange: vi.fn(),
        onCustomOpChange: vi.fn(),
        onCustomValueInput: vi.fn(),
        onApplyCustom: vi.fn(),
      },
    ).render();
    expect(filters.querySelector(".sum-popover--filters")).toBeTruthy();
    expect(filters.textContent).toContain("Draft");

    const group = renderGroupPopover(
      {
        query: { search: "", presetFilters: [], customDomain: "", groupBy: ["state"] },
        groupPresets: [{ name: "by_state", string: "Status", groupBy: "state" }],
        groupByFields: [{ name: "user_id", string: "User", type: "many2one" }],
      },
      { onToggleGroupBy: vi.fn() },
    ).render();
    expect(group.querySelector(".sum-popover--group")).toBeTruthy();
    expect(group.textContent).toContain("Status");
  });

  it("renderPopoverItem marks active entries", () => {
    const el = renderPopoverItem("Active", true, () => undefined).render();
    expect(el.querySelector(".sum-popover-item--active")).toBeTruthy();
    expect(el.textContent).toContain("✓");
  });

  it("renders favorites and actions popovers", () => {
    const favorites = renderFavoritesPopover(
      {
        favorites: [{ id: 1, name: "Mine", isShared: false }],
        saveName: "New",
        saveShared: true,
        savingFavorite: false,
      },
      {
        onApplyFavorite: vi.fn(),
        onDeleteFavorite: vi.fn(),
        onSaveNameInput: vi.fn(),
        onSaveSharedChange: vi.fn(),
        onSaveFavorite: vi.fn(),
      },
    ).render();
    expect(favorites.querySelector(".sum-popover--favorites")).toBeTruthy();
    expect(favorites.textContent).toContain("Mine");

    const actions = renderReportsPopover(basePayload(), "name").render();
    expect(actions.querySelector(".sum-popover--actions")).toBeTruthy();
    expect(actions.textContent).toContain("Import / export");
    expect(actions.textContent).toContain("Export CSV");
  });

  it("reportsToolbarVisible follows arch.report flags", () => {
    expect(reportsToolbarVisible(undefined)).toBe(false);
    expect(reportsToolbarVisible({ download: true, upload: false, pdfSizes: "", bulkModes: "" })).toBe(true);
    expect(reportsToolbarVisible({ download: false, upload: true, pdfSizes: "", bulkModes: "" })).toBe(true);
    expect(reportsToolbarVisible({ download: false, upload: false, pdfSizes: "", bulkModes: "" })).toBe(false);
  });

  it("renderReportsAnchor shows toolbar when report enabled but menu entries empty", () => {
    const payload = {
      ...basePayload(),
      arch: {
        ...basePayload().arch,
        report: { download: false, upload: true, formats: "", pdfSizes: "", bulkModes: "" },
      },
    };
    const el = renderReportsAnchor({
      open: true,
      onToggle: vi.fn(),
      payload,
      fieldsCsv: "",
      variant: "collection",
    }).render();
    expect(el.querySelector(".sum-control-bar-reports-btn")).toBeTruthy();
    expect(el.textContent).toContain("No export actions available");
  });

  it("renderReportsAnchor toggles collection reports popover", () => {
    const closed = renderReportsAnchor({
      open: false,
      onToggle: vi.fn(),
      payload: basePayload(),
      fieldsCsv: "name",
      variant: "collection",
    }).render();
    expect(closed.querySelector(".sum-popover")).toBeFalsy();
    expect(closed.querySelector(".sum-control-bar-reports-btn")).toBeTruthy();

    const open = renderReportsAnchor({
      open: true,
      onToggle: vi.fn(),
      payload: basePayload(),
      fieldsCsv: "name",
      variant: "collection",
    }).render();
    expect(open.querySelector(".sum-popover--actions")).toBeTruthy();
    expect(open.querySelector(".sum-control-bar-reports-btn--active")).toBeTruthy();
  });

  it("renderReportsAnchor renders form variant with extra popover class", () => {
    const el = renderReportsAnchor({
      open: true,
      onToggle: vi.fn(),
      payload: basePayload(),
      fieldsCsv: "name",
      variant: "form",
    }).render();
    expect(el.querySelector(".sum-form-reports-anchor")).toBeTruthy();
    expect(el.querySelector(".sum-form-reports-popover")).toBeTruthy();
  });

  it("filter popover uses readable control classes", () => {
    const filters = renderFiltersPopover(
      {
        query: { search: "", presetFilters: [], customDomain: "", groupBy: [] },
        customField: "name",
        customOp: "=",
        customValue: "",
        domainPresets: [],
        filterFields: [{ name: "name", string: "Name", type: "char" }],
      },
      {
        onTogglePreset: vi.fn(),
        onCustomFieldChange: vi.fn(),
        onCustomOpChange: vi.fn(),
        onCustomValueInput: vi.fn(),
        onApplyCustom: vi.fn(),
      },
    ).render();
    expect(filters.querySelector(".sum-popover-input")).toBeTruthy();
    expect(filters.querySelector(".sum-popover-custom .sum-btn")).toBeTruthy();
  });
});
