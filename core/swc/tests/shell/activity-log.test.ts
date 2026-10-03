import { beforeEach, describe, expect, it, vi } from "vitest";
import { initActivityLog } from "../../src/shell/activity-log.js";
import { ACTIVITY_CONTEXT } from "../../src/shell/activity-context.js";
import type { SwcEnv } from "../../src/runtime/env.js";
import { flushScheduledRenders } from "../../src/runtime/scheduler.js";

function mockEnv(getJSON: () => Promise<unknown>): SwcEnv {
  const handlers: Array<(payload: unknown) => void> = [];
  return {
    bootstrap: { swcApiBase: "/web/swc" },
    services: {
      http: { getJSON: vi.fn(getJSON) },
      bus: {
        subscribe: (topic: string, fn: (payload: unknown) => void) => {
          if (topic === ACTIVITY_CONTEXT) handlers.push(fn);
          return () => {};
        },
        emit: (_topic: string, payload: unknown) => {
          for (const fn of handlers) fn(payload);
        },
      },
    },
  } as unknown as SwcEnv;
}

describe("activity-log", () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <button type="button" data-sum-activity-tab="log" hidden>Log</button>
      <div id="sum-activity-pane-log" hidden></div>
    `;
  });

  it("shows Log tab and empty state when recordLogEligible", async () => {
    const env = mockEnv(async () => ({ items: [], hasLog: true }));
    initActivityLog(env);
    env.services.bus.emit(ACTIVITY_CONTEXT, {
      model: "crm.lead",
      recordId: 5,
      recordMessagesEligible: false,
      recordLogEligible: true,
    });
    await vi.waitFor(() => {
      flushScheduledRenders();
      const tab = document.querySelector('[data-sum-activity-tab="log"]') as HTMLElement;
      expect(tab.hidden).toBe(false);
    });
    const pane = document.getElementById("sum-activity-pane-log") as HTMLElement;
    expect(pane.hidden).toBe(false);
    expect(pane.textContent).toContain("No changes logged yet");
  });
});
