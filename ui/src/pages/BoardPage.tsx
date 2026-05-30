import { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import {
  SPEC_COLUMNS,
  columnForStatus,
  type Spec,
} from "../api/client";
import { useSpecs } from "../hooks/useSpecs";

type ViewMode = "kanban" | "table";

const VIEW_KEY = "flexspec.boardView";
const BOARD_DEFAULT_KEY = "flexspec.boardDefault";

function loadView(): ViewMode {
  const v = localStorage.getItem(VIEW_KEY);
  if (v === "table" || v === "kanban") return v;
  const fallback = localStorage.getItem(BOARD_DEFAULT_KEY);
  return fallback === "table" ? "table" : "kanban";
}

export function BoardPage() {
  const { specs, loading, error } = useSpecs();
  const [view, setView] = useState<ViewMode>(loadView);

  const setViewMode = (mode: ViewMode) => {
    setView(mode);
    localStorage.setItem(VIEW_KEY, mode);
  };

  const grouped = useMemo(() => {
    const map = new Map<string, Spec[]>();
    map.set("unassigned", []);
    for (const col of SPEC_COLUMNS) map.set(col, []);
    for (const s of specs) {
      const col = columnForStatus(s.status);
      map.get(col)!.push(s);
    }
    return map;
  }, [specs]);

  if (loading) return <p>Loading specs…</p>;
  if (error) return <p style={{ color: "#f87171" }}>{error}</p>;

  return (
    <div>
      <div style={{ display: "flex", gap: "0.5rem", marginBottom: "1rem", alignItems: "center" }}>
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
      </div>

      {view === "kanban" ? (
        <div style={{ display: "flex", gap: "0.75rem", overflowX: "auto", alignItems: "flex-start" }}>
          {[...SPEC_COLUMNS, "unassigned" as const].map((col) => (
            <div key={col} style={{ minWidth: 220, flex: "0 0 auto" }}>
              <h3 style={{ textTransform: "capitalize", fontSize: "0.85rem", color: "var(--muted)" }}>
                {col.replace(/_/g, " ")} ({grouped.get(col)?.length ?? 0})
              </h3>
              <div style={{ display: "flex", flexDirection: "column", gap: "0.5rem" }}>
                {(grouped.get(col) ?? []).map((s) => (
                  <SpecCard key={s.dir} spec={s} />
                ))}
              </div>
            </div>
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
    <div className="card">
      <div style={{ fontSize: "0.75rem", color: "var(--muted)" }}>{spec.id}</div>
      <Link to={`/specs/${spec.dir}`} style={{ fontWeight: 600 }}>
        {spec.name || spec.dir}
      </Link>
      <p style={{ margin: "0.35rem 0 0", fontSize: "0.85rem", color: "var(--muted)" }}>
        {truncate(spec.description, 100)}
      </p>
      {spec.spec_type && (
        <span style={{ fontSize: "0.7rem", opacity: 0.8 }}>{spec.spec_type}</span>
      )}
    </div>
  );
}

function truncate(s: string, n: number) {
  if (!s) return "";
  return s.length <= n ? s : s.slice(0, n) + "…";
}
