// A small fixed palette used only as tiny category dots (never big surfaces)
// across Tasks/Eisenhower/Kanban. The mapping is deterministic (hash of the
// category name) so the same category always gets the same color, without
// needing the user to pick one.
const PALETTE = [
  "#ef4444", // red
  "#22c55e", // green
  "#3b82f6", // blue
  "#a855f7", // violet
  "#f59e0b", // orange
  "#06b6d4", // cyan
];

export function categoryColor(category: string): string {
  let hash = 0;
  for (let i = 0; i < category.length; i++) {
    hash = (hash * 31 + category.charCodeAt(i)) >>> 0;
  }
  return PALETTE[hash % PALETTE.length];
}
