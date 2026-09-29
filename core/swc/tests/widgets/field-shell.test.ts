import { describe, expect, it } from "vitest";
import { html } from "../../src/template/html.js";
import {
  fieldAutocomplete,
  fieldInputId,
  fieldLabelId,
  fieldReadonlyInput,
  renderFieldShell,
  shellOptions,
} from "../../src/widgets/field-shell.js";
import type { SwcArchField } from "../../src/types/workspace.js";
import { SwcRecord } from "../../src/model/record.js";

const field: SwcArchField = { name: "date_start", type: "date", string: "Start Date" };

describe("field-shell label association", () => {
  it("editable shell links label for to explicit control id", () => {
    const id = fieldInputId(field);
    const root = renderFieldShell(field, html`<input id=${id} type="date" />`, { labelFor: id }).render();
    const label = root.querySelector("label");
    const input = root.querySelector(`#${id}`);
    expect(label?.getAttribute("for")).toBe(id);
    expect(label?.getAttribute("id")).toBe(fieldLabelId(field));
    expect(input?.id).toBe(id);
  });

  it("readonly div shell uses span label and aria-labelledby", () => {
    const root = renderFieldShell(field, html`<div class="sum-field-value">Jan 1</div>`, {
      labelFor: false,
    }).render();
    expect(root.querySelector("label")).toBeNull();
    expect(root.querySelector("span.sum-field-label")?.getAttribute("id")).toBe(fieldLabelId(field));
    expect(root.querySelector(".sum-field-control")?.getAttribute("aria-labelledby")).toBe(fieldLabelId(field));
  });

  it("fieldReadonlyInput includes matching id and autocomplete", () => {
    const root = renderFieldShell(field, fieldReadonlyInput(field, "Jan 1"), {
      labelFor: fieldInputId(field),
    }).render();
    expect(root.querySelector("label")?.getAttribute("for")).toBe(fieldInputId(field));
    const input = root.querySelector("input");
    expect(input?.id).toBe(fieldInputId(field));
    expect(input?.getAttribute("autocomplete")).toBe("off");
  });
});

describe("field-shell debug info", () => {
  const record = new SwcRecord("test.model", 1, {});

  it("renders info button when debug mode and modelName via shellOptions", () => {
    history.replaceState({}, "", "/web?debug=1");
    sessionStorage.setItem("sum.debug.mode", "1");
    const id = fieldInputId(field);
    const root = renderFieldShell(
      field,
      html`<input id=${id} type="date" />`,
      shellOptions(record, { labelFor: id }),
    ).render();
    expect(root.querySelector(".sum-debug-field-info-btn")).not.toBeNull();
    expect(root.getAttribute("data-sum-debug-field")).toBe("test.model.date_start");
    history.replaceState({}, "", "/web");
    sessionStorage.removeItem("sum.debug.mode");
  });

  it("floating debug anchor when showLabel is false", () => {
    history.replaceState({}, "", "/web?debug=1");
    sessionStorage.setItem("sum.debug.mode", "1");
    const root = renderFieldShell(field, html`<div>toggle</div>`, shellOptions(record, { showLabel: false }))
      .render();
    expect(root.querySelector(".sum-field-widget__debug-anchor .sum-debug-field-info-btn")).not.toBeNull();
    history.replaceState({}, "", "/web");
    sessionStorage.removeItem("sum.debug.mode");
  });
});

describe("fieldAutocomplete", () => {
  it("maps common field names to HTML tokens", () => {
    expect(fieldAutocomplete({ name: "email", type: "char" })).toBe("email");
    expect(fieldAutocomplete({ name: "phone", type: "char" })).toBe("tel");
    expect(fieldAutocomplete({ name: "name", type: "char" })).toBe("organization");
    expect(fieldAutocomplete({ name: "date_start", type: "date" })).toBe("off");
  });
});
