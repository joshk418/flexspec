import { useEffect, useState } from "react";
import { fetchConfigRaw, saveConfigYAML } from "../api/client";
import { Select } from "../components/Select";

const THEME_KEY = "flexspec.theme";
const BOARD_DEFAULT_KEY = "flexspec.boardDefault";

export function SettingsPage() {
  const [theme, setTheme] = useState(() => localStorage.getItem(THEME_KEY) || "dark");
  const [boardDefault, setBoardDefault] = useState(
    () => localStorage.getItem(BOARD_DEFAULT_KEY) || "kanban"
  );
  const [yaml, setYaml] = useState("");
  const [configError, setConfigError] = useState<string | null>(null);
  const [configOk, setConfigOk] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    document.documentElement.dataset.theme = theme;
    localStorage.setItem(THEME_KEY, theme);
  }, [theme]);

  useEffect(() => {
    localStorage.setItem(BOARD_DEFAULT_KEY, boardDefault);
  }, [boardDefault]);

  useEffect(() => {
    fetchConfigRaw()
      .then(setYaml)
      .catch((e) => setConfigError(e instanceof Error ? e.message : "Load failed"))
      .finally(() => setLoading(false));
  }, []);

  const saveConfig = async () => {
    setConfigError(null);
    setConfigOk(null);
    try {
      await saveConfigYAML(yaml);
      setConfigOk("Config saved.");
    } catch (e) {
      setConfigError(e instanceof Error ? e.message : "Save failed");
    }
  };

  return (
    <div style={{ maxWidth: 720 }}>
      <h1 style={{ marginTop: 0 }}>Settings</h1>

      <section className="card" style={{ marginBottom: "1rem", display: "flex", flexDirection: "column", gap: "1rem" }}>
        <h2 style={{ marginTop: 0, marginBottom: "0.25rem", fontSize: "1rem" }}>Appearance</h2>
        
        <div>
          <span style={{ display: "block", marginBottom: "0.35rem", fontSize: "0.9rem", color: "var(--muted)" }}>Theme</span>
          <Select
            value={theme}
            onChange={setTheme}
            options={[
              { value: "dark", label: "Dark" },
              { value: "light", label: "Light" },
            ]}
          />
        </div>

        <div>
          <span style={{ display: "block", marginBottom: "0.35rem", fontSize: "0.9rem", color: "var(--muted)" }}>Default board view</span>
          <Select
            value={boardDefault}
            onChange={setBoardDefault}
            options={[
              { value: "kanban", label: "Kanban" },
              { value: "table", label: "Table" },
            ]}
          />
        </div>
      </section>

      <section className="card">
        <h2 style={{ marginTop: 0, fontSize: "1rem" }}>FlexSpec config</h2>
        <p style={{ color: "var(--muted)", fontSize: "0.9rem" }}>
          Edits <code>.flexspec/config.yaml</code>. Invalid YAML or values return an error.
        </p>
        {loading ? (
          <p>Loading config…</p>
        ) : (
          <>
            <textarea
              value={yaml}
              onChange={(e) => setYaml(e.target.value)}
              rows={12}
              style={{
                width: "100%",
                fontFamily: "ui-monospace, monospace",
                fontSize: "0.85rem",
                background: "#0a0e14",
                color: "var(--text)",
                border: "1px solid var(--border)",
                borderRadius: 6,
                padding: "0.75rem",
              }}
            />
            <div style={{ marginTop: "0.75rem", display: "flex", gap: "0.5rem" }}>
              <button type="button" className="btn" onClick={saveConfig}>
                Save config
              </button>
            </div>
            {configError && <p style={{ color: "#f87171" }}>{configError}</p>}
            {configOk && <p style={{ color: "#4ade80" }}>{configOk}</p>}
          </>
        )}
      </section>
    </div>
  );
}
