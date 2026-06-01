import { useEffect, useState } from "react";
import { fetchConfig, saveConfig, type ProjectConfig } from "../api/client";
import { Select } from "../components/Select";

const THEME_KEY = "flexspec.theme";
const BOARD_DEFAULT_KEY = "flexspec.boardDefault";

export function SettingsPage() {
  const [theme, setTheme] = useState(() => localStorage.getItem(THEME_KEY) || "dark");
  const [boardDefault, setBoardDefault] = useState(
    () => localStorage.getItem(BOARD_DEFAULT_KEY) || "kanban"
  );
  const [config, setConfig] = useState<ProjectConfig | null>(null);
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
    fetchConfig()
      .then(setConfig)
      .catch((e) => setConfigError(e instanceof Error ? e.message : "Load failed"))
      .finally(() => setLoading(false));
  }, []);

  const updateConfig = (patch: Partial<ProjectConfig>) => {
    setConfig((prev) => (prev ? { ...prev, ...patch } : prev));
  };

  const handleSaveConfig = async () => {
    if (!config) return;
    setConfigError(null);
    setConfigOk(null);
    try {
      const saved = await saveConfig(config);
      setConfig(saved);
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
          Updates <code>.flexspec/config.yaml</code> through structured fields.
        </p>
        {loading ? (
          <p>Loading config…</p>
        ) : config ? (
          <>
            <table className="config-table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>specs_dir</td>
                  <td>
                    <input
                      type="text"
                      className="config-input"
                      value={config.specs_dir}
                      onChange={(e) => updateConfig({ specs_dir: e.target.value })}
                    />
                  </td>
                </tr>
                <tr>
                  <td>always_one_shot</td>
                  <td>
                    <Select
                      value={config.always_one_shot ? "true" : "false"}
                      onChange={(v) => updateConfig({ always_one_shot: v === "true" })}
                      options={[
                        { value: "false", label: "false" },
                        { value: "true", label: "true" },
                      ]}
                    />
                  </td>
                </tr>
                <tr>
                  <td>spec_template</td>
                  <td>
                    <Select
                      value={config.spec_template}
                      onChange={(v) => updateConfig({ spec_template: v })}
                      options={[
                        { value: "", label: "Infer" },
                        { value: "simple", label: "Simple" },
                        { value: "expanded", label: "Expanded" },
                      ]}
                    />
                  </td>
                </tr>
              </tbody>
            </table>
            <div style={{ marginTop: "0.75rem", display: "flex", gap: "0.5rem" }}>
              <button type="button" className="btn" onClick={handleSaveConfig}>
                Save config
              </button>
            </div>
            {configError && <p style={{ color: "#f87171" }}>{configError}</p>}
            {configOk && <p style={{ color: "#4ade80" }}>{configOk}</p>}
          </>
        ) : (
          configError && <p style={{ color: "#f87171" }}>{configError}</p>
        )}
      </section>
    </div>
  );
}
