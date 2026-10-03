import type { SwcWorkspacePayload } from "../types/workspace.js";

export const ACTIVITY_CONTEXT = "activity.context";

export interface ActivityContextPayload {
  model: string;
  recordId: number;
  recordMessagesEligible: boolean;
  recordLogEligible: boolean;
}

export function activityContextFromWorkspace(payload: SwcWorkspacePayload): ActivityContextPayload {
  return {
    model: payload.model,
    recordId: payload.recordId,
    recordMessagesEligible: payload.recordMessagesEligible === true,
    recordLogEligible: payload.recordLogEligible === true,
  };
}

export function emptyActivityContext(): ActivityContextPayload {
  return {
    model: "",
    recordId: 0,
    recordMessagesEligible: false,
    recordLogEligible: false,
  };
}

export function messagesEmptyStateHTML(): string {
  return `<div class="sum-msg-empty" role="status">
    <div class="sum-msg-empty-visual" aria-hidden="true">
      <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.25">
        <path d="M21 15a2 2 0 01-2 2H7l-4 4V5a2 2 0 012-2h14a2 2 0 012 2z" />
      </svg>
    </div>
    <p class="sum-msg-empty-title">Messages</p>
    <p class="sum-msg-empty-hint">Search for an internal user to start a direct conversation.</p>
  </div>`;
}
