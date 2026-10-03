const PREFIX = "sum-im-draft:";

export function readImDraft(peerId: number): string {
  if (peerId <= 0) return "";
  try {
    return localStorage.getItem(PREFIX + String(peerId)) ?? "";
  } catch {
    return "";
  }
}

export function writeImDraft(peerId: number, text: string): void {
  if (peerId <= 0) return;
  try {
    const key = PREFIX + String(peerId);
    const trimmed = text.trim();
    if (!trimmed) localStorage.removeItem(key);
    else localStorage.setItem(key, text);
  } catch {
    /* ignore */
  }
}

export function clearImDraft(peerId: number): void {
  writeImDraft(peerId, "");
}
