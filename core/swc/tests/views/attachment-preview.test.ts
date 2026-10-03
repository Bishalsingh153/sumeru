import { describe, expect, it } from "vitest";
import { attachmentContentInlineUrl } from "../../src/views/shared/AttachmentPreview.js";

describe("attachmentContentInlineUrl", () => {
  it("appends download=0 for inline preview", () => {
    expect(attachmentContentInlineUrl("/web/content/5")).toBe("/web/content/5?download=0");
  });

  it("does not duplicate download param", () => {
    expect(attachmentContentInlineUrl("/web/content/5?download=0")).toBe("/web/content/5?download=0");
  });
});
