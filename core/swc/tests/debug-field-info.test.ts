import { describe, expect, it } from "vitest";
import {
  buildFieldDebugInfo,
  formatDebugContext,
  formatDebugDomain,
  formatDebugFieldValue,
  isSensitiveDebugField,
  modifierSummary,
  parseDebugFieldKey,
} from "../src/devtools/debug-field-info.js";
import type { SwcArchField } from "../src/types/workspace.js";

describe("debug field info", () => {
  it("parseDebugFieldKey splits model and field", () => {
    expect(parseDebugFieldKey("res.partner.email")).toEqual({
      model: "res.partner",
      fieldName: "email",
    });
  });

  it("redacts sensitive field values", () => {
    expect(isSensitiveDebugField("password")).toBe(true);
    expect(isSensitiveDebugField("api_key")).toBe(true);
    expect(formatDebugFieldValue("password", "secret123")).toBe("— (redacted)");
    expect(formatDebugFieldValue("name", "Acme")).toBe("Acme");
  });

  it("formatDebugDomain parses arch option domain", () => {
    const field: SwcArchField = {
      name: "country_id",
      options: { domain: '[["id","=",1]]' },
    };
    expect(formatDebugDomain(field)).toContain("id");
  });

  it("formatDebugContext empty object", () => {
    expect(formatDebugContext({})).toBe("{}");
  });

  it("buildFieldDebugInfo merges arch and ACL", () => {
    const arch: SwcArchField = {
      name: "state",
      string: "Status",
      type: "selection",
      widget: "statusbar",
      readonly_expr: "state == 'done'",
    };
    const info = buildFieldDebugInfo({
      fieldKey: "sale.order.state",
      archField: arch,
      recordData: { state: "draft" },
      aclField: { write_denied: true },
      userId: 1,
      companyId: 1,
    });
    expect(info.label).toBe("Status");
    expect(info.readonlyExpr).toBe("state == 'done'");
    expect(info.aclWriteDenied).toBe(true);
    expect(info.valueDisplay).toBe("draft");
    expect(modifierSummary(info)).toBe("none");
  });
});
