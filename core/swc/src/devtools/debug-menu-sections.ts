import type { DebugWorkspaceContext } from "./debug.js";

export interface DebugMenuSection {
  id: string;
  title: string;
  items: string[];
}

/** Section titles shown when developer mode is active (for tests and menu builder). */
export function buildDebugMenuSections(
  debugOn: boolean,
  workspace?: DebugWorkspaceContext,
): DebugMenuSection[] {
  if (!debugOn) return [];
  const sections: DebugMenuSection[] = [
    {
      id: "record",
      title: "Record",
      items: ["Metadata", "Data", "Messages", "Attachments", "Set default values"],
    },
    {
      id: "ui",
      title: "User interface",
      items: [
        workspace?.model ? `Model: ${workspace.model}` : "Model: —",
        workspace?.viewType ? `View: ${workspace.viewType}` : "View: —",
        "Fields",
        "Filters",
        "View arch",
      ],
    },
    { id: "security", title: "Security", items: ["Access rights"] },
    {
      id: "tools",
      title: "Tools",
      items: ["Regenerate assets", "Field inspector", "Open debug drawer", "Open SWC Vision", "Open metrics", "Reload page"],
    },
  ];
  return sections;
}
