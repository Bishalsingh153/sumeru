export function userInitials(name: string, login = ""): string {
  const nm = name.trim();
  if (nm) {
    const parts = nm.split(/\s+/).filter(Boolean);
    if (parts.length >= 2) return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
    return nm.slice(0, 2).toUpperCase();
  }
  const lg = login.trim();
  if (lg.length >= 2) return lg.slice(0, 2).toUpperCase();
  return lg.slice(0, 1).toUpperCase() || "?";
}
