import { describe, expect, it } from "vitest";
import { buildDebugMenuSections } from "../src/devtools/debug-menu-sections.js";

describe("debug menu sections", () => {
  it("returns empty when debug off", () => {
    expect(buildDebugMenuSections(false)).toEqual([]);
  });

  it("includes model and view readouts when workspace set", () => {
    const sections = buildDebugMenuSections(true, {
      model: "core.partner",
      recordId: 1,
      viewType: "form",
    });
    const ui = sections.find((s) => s.id === "ui");
    expect(ui?.items.some((i) => i.includes("core.partner"))).toBe(true);
    expect(ui?.items.some((i) => i.includes("form"))).toBe(true);
  });
});
