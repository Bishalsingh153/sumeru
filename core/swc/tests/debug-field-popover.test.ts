import { describe, expect, it, beforeEach, afterEach } from "vitest";
import { fieldDebugPopoverLines, buildFieldDebugInfo } from "../src/devtools/debug-field-info.js";
import { hideFieldDebugPopover, showFieldDebugPopover } from "../src/devtools/debug-field-popover.js";

describe("debug field popover", () => {
  beforeEach(() => {
    document.body.innerHTML = "";
  });

  afterEach(() => {
    hideFieldDebugPopover();
    document.getElementById("sum-debug-field-popover")?.remove();
  });

  it("fieldDebugPopoverLines includes domain and relation", () => {
    const info = buildFieldDebugInfo({
      fieldKey: "core.partner.name",
      archField: { name: "name", string: "Name", type: "char", relation: "—" },
    });
    const lines = fieldDebugPopoverLines(info);
    expect(lines.some((l) => l.startsWith("Field:"))).toBe(true);
    expect(lines.some((l) => l.startsWith("Domain:"))).toBe(true);
  });

  it("showFieldDebugPopover renders popover content", () => {
    const wrap = document.createElement("div");
    wrap.setAttribute("data-sum-debug-field", "core.partner.email");
    wrap.setAttribute("data-sum-debug-label", "Email");
    wrap.setAttribute("data-sum-debug-type", "char");
    const btn = document.createElement("button");
    btn.className = "sum-debug-field-info-btn";
    wrap.appendChild(btn);
    document.body.appendChild(wrap);
    showFieldDebugPopover(btn);
    const pop = document.getElementById("sum-debug-field-popover");
    expect(pop?.hidden).toBe(false);
    expect(pop?.textContent).toContain("email");
  });
});
