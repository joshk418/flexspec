import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  SPEC_COLUMNS,
  UNASSIGNED_COLUMN,
  columnForStatus,
  type Spec,
} from "../api/client";
import { useSpecs } from "../hooks/useSpecs";

type ViewMode = "kanban" | "table";

const VIEW_KEY = "flexspec.boardView";
const BOARD_DEFAULT_KEY = "flexspec.boardDefault";
const COLUMNS_KEY = "flexspec.boardColumns";

const ALL_COLUMNS = [...SPEC_COLUMNS, UNASSIGNED_COLUMN] as const;

function loadView(): ViewMode {
  const v = localStorage.getItem(VIEW_KEY);
  if (v === "table" || v === "kanban") return v;
  const fallback = localStorage.getItem(BOARD_DEFAULT_KEY);
  return fallback === "table" ? "table" : "kanban";
}

function loadVisibleColumns(): string[] {
  try {
    const raw = localStorage.getItem(COLUMNS_KEY);
    if (raw) {
      const parsed = JSON.parse(raw);
      if (Array.isArray(parsed)) {
        const valid = parsed.filter((c) => (ALL_COLUMNS as readonly string[]).includes(c));
        if (valid.length > 0) return valid;
      }
    }
  } catch {
    /* fall through to default */
  }
  return [...ALL_COLUMNS];
}

export function BoardPage() {
  const { specs, loading, error } = useSpecs();
  const [view, setView] = useState<ViewMode>(loadView);
  const [visibleColumns, setVisibleColumns] = useState<string[]>(loadVisibleColumns);

  const setViewMode = (mode: ViewMode) => {
    setView(mode);
    localStorage.setItem(VIEW_KEY, mode);
  };

  const toggleColumn = (col: string) => {
    setVisibleColumns((prev) => {
      const next = prev.includes(col) ? prev.filter((c) => c !== col) : [...prev, col];
      localStorage.setItem(COLUMNS_KEY, JSON.stringify(next));
      return next;
    });
  };

  const grouped = useMemo(() => {
    const map = new Map<string, Spec[]>();
    for (const col of ALL_COLUMNS) map.set(col, []);
    for (const s of specs) {
      map.get(columnForStatus(s.status))!.push(s);
    }
    return map;
  }, [specs]);

  // Preserve canonical column order; render only columns the user chose to show.
  const renderedColumns = useMemo(
    () => ALL_COLUMNS.filter((col) => visibleColumns.includes(col)),
    [visibleColumns],
  );

  if (loading) return <p>Loading specs…</p>;
  if (error) return <p style={{ color: "#f87171" }}>{error}</p>;

  return (
    <div>
      <div className="board-toolbar">
        <strong>Board</strong>
        <button
          type="button"
          className={`btn secondary ${view === "kanban" ? "active" : ""}`}
          onClick={() => setViewMode("kanban")}
        >
          Kanban
        </button>
        <button
          type="button"
          className={`btn secondary ${view === "table" ? "active" : ""}`}
          onClick={() => setViewMode("table")}
        >
          Table
        </button>

        {view === "kanban" && (
          <details className="board-columns-menu">
            <summary className="btn secondary">Columns</summary>
            <div className="board-columns-panel">
              {ALL_COLUMNS.map((col) => (
                <label key={col} className="board-columns-option">
                  <input
                    type="checkbox"
                    checked={visibleColumns.includes(col)}
                    onChange={() => toggleColumn(col)}
                  />
                  <span>{col.replace(/_/g, " ")}</span>
                </label>
              ))}
            </div>
          </details>
        )}
      </div>

      {view === "kanban" ? (
        <div
          className="board-kanban"
          style={{ gridTemplateColumns: `repeat(${Math.max(renderedColumns.length, 1)}, minmax(0, 1fr))` }}
        >
          {renderedColumns.map((col) => (
            <section key={col} className="board-column">
              <h3 className="board-column-title">
                {col.replace(/_/g, " ")} ({grouped.get(col)?.length ?? 0})
              </h3>
              <div className="board-column-cards">
                {(grouped.get(col) ?? []).map((s) => (
                  <SpecCard key={s.dir} spec={s} />
                ))}
              </div>
            </section>
          ))}
        </div>
      ) : (
        <table style={{ width: "100%", borderCollapse: "collapse" }}>
          <thead>
            <tr style={{ textAlign: "left", color: "var(--muted)" }}>
              <th>ID</th>
              <th>Name</th>
              <th>Status</th>
              <th>Type</th>
              <th>Description</th>
            </tr>
          </thead>
          <tbody>
            {specs.map((s) => (
              <tr key={s.dir} style={{ borderTop: "1px solid var(--border)" }}>
                <td>{s.id}</td>
                <td>
                  <Link to={`/specs/${s.dir}`}>{s.name || s.dir}</Link>
                </td>
                <td>{s.status || "—"}</td>
                <td>{s.spec_type || "—"}</td>
                <td style={{ color: "var(--muted)", maxWidth: 400 }}>{truncate(s.description, 80)}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );
}

function SpecCard({ spec }: { spec: Spec }) {
  return (
    <div className="card board-card">
      <div className="board-card-id">{spec.id}</div>
      <Link to={`/specs/${spec.dir}`} className="board-card-title">
        {spec.name || spec.dir}
      </Link>
      <p className="board-card-desc">{truncate(spec.description, 100)}</p>
      {spec.spec_type && <span className="board-card-type">{spec.spec_type}</span>}
    </div>
  );
}

function truncate(s: string, n: number) {
  if (!s) return "";
  return s.length <= n ? s : s.slice(0, n) + "…";
}
