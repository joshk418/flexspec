// Canonical spec lifecycle statuses in board order.
// Keep in sync with internal/spec/status.go (SpecStatuses).
export const SPEC_COLUMNS = [
  "draft",
  "planned",
  "in_progress",
  "in_review",
  "complete",
] as const;

export type SpecColumn = (typeof SPEC_COLUMNS)[number];

export const UNASSIGNED_COLUMN = "unassigned";

// Legacy/renamed statuses mapped to their current value.
// Keep in sync with internal/spec/status.go (legacyStatusAliases).
const LEGACY_STATUS_ALIASES: Record<string, string> = {
  refined: "planned",
  initial: "draft",
};

export function normalizeSpecStatus(status: string): string {
  const s = status.trim().toLowerCase();
  return LEGACY_STATUS_ALIASES[s] ?? s;
}

// columnForStatus normalizes legacy statuses first, then maps anything unknown
// to the Unassigned column.
export function columnForStatus(status: string): string {
  const s = normalizeSpecStatus(status);
  return (SPEC_COLUMNS as readonly string[]).includes(s) ? s : UNASSIGNED_COLUMN;
}
