import { describe, expect, it, vi } from "vitest";
import { InternalMessagesPanel } from "../../src/shell/internal-messages-panel.js";
import { SwcError } from "../../src/runtime/error.js";
import type { SwcEnv } from "../../src/runtime/env.js";
import { flushScheduledRenders } from "../../src/runtime/scheduler.js";

function panelEnv(fetchImpl: (url: string) => Promise<unknown>): SwcEnv {
  return {
    bootstrap: { swcApiBase: "/web/swc", internalMessagesEnabled: true },
    services: {
      http: {
        getJSON: vi.fn((url: string) => fetchImpl(url)),
        postForm: vi.fn(),
        postMultipart: vi.fn(),
      },
      bus: { subscribe: () => () => {} },
    },
  } as unknown as SwcEnv;
}

describe("InternalMessagesPanel", () => {
  it("keeps search UI when conversations fail but directory loads", async () => {
    const env = panelEnv(async (url) => {
      if (url.includes("/direct/status")) {
        return { enabled: true, unread: 0 };
      }
      if (url.includes("/direct/conversations")) {
        throw new SwcError("GET failed: 500", "http_get");
      }
      if (url.includes("/direct/users")) {
        return { users: [{ id: 2, name: "Alex", login: "alex" }] };
      }
      return {};
    });
    const panel = new InternalMessagesPanel({}, env);
    panel.callSetup();
    const mount = document.createElement("div");
    document.body.appendChild(mount);
    mount.replaceChildren(panel.render());
    flushScheduledRenders();
    await vi.waitFor(() => {
      flushScheduledRenders();
      expect(mount.querySelector(".sum-msg-search")).toBeTruthy();
      expect(mount.textContent).toContain("Alex");
    });
    expect(mount.querySelector(".sum-msg-empty")).toBeFalsy();
  });

  it("shows company link when IM disabled via status", async () => {
    const env = panelEnv(async (url) => {
      if (url.includes("/direct/status")) {
        return { enabled: false, companyFormHref: "/web?model=core.company&id=1" };
      }
      return {};
    });
    env.bootstrap.internalMessagesEnabled = false;
    const panel = new InternalMessagesPanel({}, env);
    panel.callSetup();
    const mount = document.createElement("div");
    document.body.appendChild(mount);
    mount.replaceChildren(panel.render());
    flushScheduledRenders();
    await vi.waitFor(() => {
      flushScheduledRenders();
      expect(mount.textContent).toContain("company settings");
    });
  });

  it("fuzzy search surfaces John for jn", async () => {
    const env = panelEnv(async (url) => {
      if (url.includes("/direct/status")) return { enabled: true };
      if (url.includes("/direct/conversations")) return { conversations: [] };
      if (url.includes("/direct/users")) {
        return {
          users: [
            { id: 1, name: "Jane Doe", login: "jane" },
            { id: 2, name: "John Smith", login: "john" },
          ],
        };
      }
      return {};
    });
    const panel = new InternalMessagesPanel({}, env);
    panel.callSetup();
    const mount = document.createElement("div");
    document.body.appendChild(mount);
    mount.replaceChildren(panel.render());
    await vi.waitFor(() => {
      flushScheduledRenders();
      expect(mount.querySelector(".sum-msg-search")).toBeTruthy();
    });
    const input = mount.querySelector(".sum-msg-search") as HTMLInputElement;
    input.value = "jn";
    input.dispatchEvent(new Event("input", { bubbles: true }));
    flushScheduledRenders();
    await vi.waitFor(() => {
      flushScheduledRenders();
      expect(mount.textContent).toContain("John");
    });
  });
});
