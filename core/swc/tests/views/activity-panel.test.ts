import { beforeEach, describe, expect, it } from "vitest";
import { applyActivityHidden, setActivityLogTabVisible } from "../../src/shell/activity-panel.js";

describe("activity-panel", () => {
  beforeEach(() => {
    document.body.innerHTML = `
      <div id="sum-shell" class="sum-shell">
        <aside id="sum-activity-panel" class="sum-activity-panel">
          <button type="button" class="sum-activity-tab is-active" data-sum-activity-tab="messages">Messages</button>
          <button type="button" class="sum-activity-tab" data-sum-activity-tab="log" hidden>Log</button>
          <div id="sum-activity-pane-messages"></div>
          <div id="sum-activity-pane-log" hidden></div>
        </aside>
        <button type="button" id="sum-activity-reveal" hidden></button>
        <button type="button" id="sum-activity-toggle" aria-pressed="true"></button>
      </div>
    `;
  });

  it("applyActivityHidden toggles shell class and reveal button", () => {
    const shell = document.getElementById("sum-shell") as HTMLElement;
    applyActivityHidden(shell, true);
    expect(shell.classList.contains("sum-shell--activity-hidden")).toBe(true);
    expect((document.getElementById("sum-activity-reveal") as HTMLButtonElement).hidden).toBe(false);

    applyActivityHidden(shell, false);
    expect(shell.classList.contains("sum-shell--activity-hidden")).toBe(false);
    expect((document.getElementById("sum-activity-reveal") as HTMLButtonElement).hidden).toBe(true);
  });

  it("setActivityLogTabVisible toggles log tab and pane", () => {
    const logTab = document.querySelector('[data-sum-activity-tab="log"]') as HTMLElement;
    setActivityLogTabVisible(true);
    expect(logTab.hidden).toBe(false);
    setActivityLogTabVisible(false);
    expect(logTab.hidden).toBe(true);
    expect((document.getElementById("sum-activity-pane-log") as HTMLElement).innerHTML).toBe("");
  });
});
