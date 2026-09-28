import { describe, expect, it, beforeEach, afterEach } from "vitest";
import {
  applyDebugBodyClasses,
  getDebugMode,
  isFieldInspectorEnabled,
  parseDebugParam,
  selectDebugFieldFromElement,
  toggleFieldInspector,
} from "../src/devtools/debug.js";

describe("debug manager", () => {
  beforeEach(() => {
    sessionStorage.clear();
    document.body.className = "";
    document.body.innerHTML = "";
    history.replaceState({}, "", "/web");
  });

  afterEach(() => {
    sessionStorage.clear();
    document.body.className = "";
  });

  it("parseDebugParam accepts common truthy values", () => {
    expect(parseDebugParam("1")).toBe("1");
    expect(parseDebugParam("true")).toBe("1");
    expect(parseDebugParam("assets")).toBe("assets");
    expect(parseDebugParam("no")).toBe("off");
  });

  it("applyDebugBodyClasses reflects mode and field inspector", () => {
    sessionStorage.setItem("sum.debug.mode", "1");
    history.replaceState({}, "", "/web?debug=1");
    toggleFieldInspector();
    applyDebugBodyClasses();
    expect(document.body.classList.contains("sum-debug--active")).toBe(true);
    expect(document.body.classList.contains("sum-debug--field-inspector")).toBe(true);
  });

  it("selectDebugFieldFromElement marks selected field", () => {
    const el = document.createElement("div");
    el.setAttribute("data-sum-debug-field", "res.partner.name");
    document.body.appendChild(el);
    selectDebugFieldFromElement(el);
    expect(el.classList.contains("sum-debug-field--selected")).toBe(true);
  });

  it("getDebugMode prefers URL param over session storage", () => {
    sessionStorage.setItem("sum.debug.mode", "assets");
    history.replaceState({}, "", "/web?debug=1");
    expect(getDebugMode()).toBe("1");
  });

  it("toggleFieldInspector flips session flag", () => {
    expect(isFieldInspectorEnabled()).toBe(false);
    toggleFieldInspector();
    expect(isFieldInspectorEnabled()).toBe(true);
    toggleFieldInspector();
    expect(isFieldInspectorEnabled()).toBe(false);
  });
});
