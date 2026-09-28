import { describe, expect, it, beforeEach, afterEach } from "vitest";
import { resetDebugHoverForTests, setDebugHoverTarget } from "../src/devtools/debug-hover.js";

describe("debug hover trace", () => {
  beforeEach(() => {
    document.body.innerHTML = "";
    resetDebugHoverForTests();
  });

  afterEach(() => {
    resetDebugHoverForTests();
    document.getElementById("sum-debug-hover-card")?.remove();
  });

  it("setDebugHoverTarget toggles hover class and shows card", () => {
    const field = document.createElement("div");
    field.className = "sum-field-widget";
    field.setAttribute("data-sum-debug-field", "res.partner.name");
    field.setAttribute("data-sum-debug-label", "Name");
    field.setAttribute("data-sum-debug-type", "char");
    document.body.appendChild(field);

    setDebugHoverTarget(field, 40, 40);
    expect(field.classList.contains("sum-debug-field--hover")).toBe(true);
    const card = document.getElementById("sum-debug-hover-card");
    expect(card?.hidden).toBe(false);
    expect(card?.textContent).toContain("name");

    setDebugHoverTarget(null, 0, 0);
    expect(field.classList.contains("sum-debug-field--hover")).toBe(false);
    expect(card?.hidden).toBe(true);
  });
});
