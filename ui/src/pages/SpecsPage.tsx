import { useCallback, useEffect, useState } from "react";
import { Link, useParams } from "react-router-dom";
import {
  fetchSpecDetail,
  subscribeEvents,
  type SpecDetail,
} from "../api/client";
import { SpecMarkdown } from "../components/SpecMarkdown";
import { useSpecs } from "../hooks/useSpecs";

export function SpecsPage() {
  const { specs, loading, error } = useSpecs();
  const { dir } = useParams<{ dir?: string }>();

  if (loading) return <p>Loading…</p>;
  if (error) return <p style={{ color: "#f87171" }}>{error}</p>;

  return (
    <div style={{ display: "grid", gridTemplateColumns: "240px 1fr", gap: "1rem" }}>
      <aside>
        <h2 style={{ marginTop: 0 }}>Specs</h2>
        <ul style={{ listStyle: "none", padding: 0, margin: 0 }}>
          {specs.map((s) => (
            <li key={s.dir} style={{ marginBottom: "0.35rem" }}>
              <Link
                to={`/specs/${s.dir}`}
                style={{ fontWeight: dir === s.dir ? 700 : 400 }}
              >
                {s.name || s.dir}
              </Link>
              <div style={{ fontSize: "0.75rem", color: "var(--muted)" }}>{s.status}</div>
            </li>
          ))}
        </ul>
      </aside>
      <section>
        {dir ? <SpecDetailView dir={dir} /> : <p style={{ color: "var(--muted)" }}>Select a spec</p>}
      </section>
    </div>
  );
}

function SpecDetailView({ dir }: { dir: string }) {
  const [detail, setDetail] = useState<SpecDetail | null>(null);
  const [error, setError] = useState<string | null>(null);

  const load = useCallback(async () => {
    try {
      setError(null);
      setDetail(await fetchSpecDetail(dir));
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load spec");
    }
  }, [dir]);

  useEffect(() => {
    load();
    return subscribeEvents(load);
  }, [load]);

  if (error) return <p style={{ color: "#f87171" }}>{error}</p>;
  if (!detail) return <p>Loading spec…</p>;

  return (
    <div>
      <h1 style={{ marginTop: 0 }}>{detail.name || detail.dir}</h1>
      <p style={{ color: "var(--muted)" }}>
        {detail.status} · {detail.spec_type}
      </p>
      <SpecMarkdown content={detail.markdown} />
      {detail.tasks && detail.tasks.length > 0 && (
        <div style={{ marginTop: "1.5rem" }}>
          <h2>Tasks</h2>
          {detail.tasks.map((t) => (
            <TaskAccordion key={t.file} task={t} />
          ))}
        </div>
      )}
    </div>
  );
}

function TaskAccordion({ task }: { task: SpecDetail["tasks"][0] }) {
  const [open, setOpen] = useState(false);
  return (
    <div className="card" style={{ marginBottom: "0.5rem" }}>
      <button
        type="button"
        onClick={() => setOpen(!open)}
        style={{
          width: "100%",
          textAlign: "left",
          background: "none",
          border: "none",
          color: "var(--text)",
          padding: 0,
        }}
      >
        <strong>{task.id || task.file}</strong> — {task.name || "Task"}{" "}
        <span style={{ color: "var(--muted)" }}>({task.status})</span>
        <span style={{ float: "right" }}>{open ? "▼" : "▶"}</span>
      </button>
      {open && (
        <div style={{ marginTop: "0.75rem" }}>
          <SpecMarkdown content={task.markdown} />
        </div>
      )}
    </div>
  );
}
