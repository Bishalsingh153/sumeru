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
