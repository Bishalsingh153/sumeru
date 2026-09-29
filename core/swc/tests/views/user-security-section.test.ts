import { describe, expect, it } from "vitest";
import {
  gatherUserSecurityFields,
  renderUserSecuritySection,
  validateNewUserSecurityFields,
} from "../../src/views/form/user-security-section.js";

describe("user-security-section", () => {
  const meta = {
    canEdit: true,
    isNew: true,
    passwordLabel: "Set initial password",
    groups: [
      { id: 1, name: "Internal", selected: true, isUserType: true },
      { id: 10, name: "Settings", selected: false, isUserType: false },
    ],
    companies: [{ id: 1, name: "Acme", selected: true }],
  };

  it("renders password fields for Account Security mode", () => {
    const root = renderUserSecuritySection(meta, "password", false).render();
    expect(root.querySelector("#password_plain")).toBeTruthy();
    expect(root.querySelector('[name="password_plain_confirm"]')).toBeTruthy();
  });

  it("gatherUserSecurityFields collects groups and password", () => {
    const host = document.createElement("div");
    host.appendChild(renderUserSecuritySection(meta, "access", false).render());
    host.appendChild(renderUserSecuritySection(meta, "password", false).render());
    const fields = gatherUserSecurityFields(host);
    expect(fields.security_groups_touched).toBe("1");
    expect(fields.security_user_type).toBe("1");
    expect(fields.password_plain).toBeUndefined();
    const pw = host.querySelector<HTMLInputElement>("#password_plain");
    if (pw) pw.value = "SecretPass1";
    const pwc = host.querySelector<HTMLInputElement>("#password_plain_confirm");
    if (pwc) pwc.value = "SecretPass1";
    const again = gatherUserSecurityFields(host);
    expect(again.password_plain).toBe("SecretPass1");
  });

  it("validateNewUserSecurityFields requires matching non-empty password", () => {
    expect(validateNewUserSecurityFields(null)).toMatch(/initial password/i);
    expect(validateNewUserSecurityFields({})).toMatch(/initial password/i);
    expect(
      validateNewUserSecurityFields({
        password_plain: "SecretPass1",
        password_plain_confirm: "other",
      }),
    ).toMatch(/do not match/i);
    expect(
      validateNewUserSecurityFields({
        password_plain: "SecretPass1",
        password_plain_confirm: "SecretPass1",
      }),
    ).toBeNull();
  });

  it("security snapshot keeps password after inputs are cleared (save rerender)", () => {
    const host = document.createElement("div");
    host.appendChild(renderUserSecuritySection(meta, "password", false).render());
    const pw = host.querySelector<HTMLInputElement>("#password_plain");
    const pwc = host.querySelector<HTMLInputElement>("#password_plain_confirm");
    if (!pw || !pwc) throw new Error("password inputs missing");
    pw.value = "SecretPass1";
    pwc.value = "SecretPass1";
    const snapshot = gatherUserSecurityFields(host);
    pw.value = "";
    pwc.value = "";
    expect(gatherUserSecurityFields(host).password_plain).toBeUndefined();
    expect(snapshot.password_plain).toBe("SecretPass1");
    expect(snapshot.password_plain_confirm).toBe("SecretPass1");
  });
});
