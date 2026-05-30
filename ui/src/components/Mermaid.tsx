import { useEffect, useRef, useState } from "react";
import mermaid from "mermaid";

// Initialize mermaid once with safe default settings
mermaid.initialize({
  startOnLoad: false,
  theme: "dark",
  securityLevel: "loose",
});

let uniqueId = 0;

export function Mermaid({ chart }: { chart: string }) {
  const [svg, setSvg] = useState<string>("");
  const [error, setError] = useState<string | null>(null);
  const elementId = useRef(`mermaid-${++uniqueId}`);

  useEffect(() => {
    // Dynamically adjust mermaid theme based on application's current dataset.theme
    const isLight = document.documentElement.dataset.theme === "light";
    mermaid.initialize({
      startOnLoad: false,
      theme: isLight ? "default" : "dark",
      securityLevel: "loose",
    });

    let active = true;

    async function renderChart() {
      if (!chart.trim()) return;
      try {
        setError(null);
        // Modern mermaid.render returns { svg, bindFunctions } asynchronously
        const { svg: renderedSvg } = await mermaid.render(elementId.current, chart);
        if (active) {
          setSvg(renderedSvg);
        }
      } catch (err) {
        console.error("Mermaid render error:", err);
        if (active) {
          setError(err instanceof Error ? err.message : String(err));
        }
      }
    }

    renderChart();

    return () => {
      active = false;
    };
  }, [chart]);

  if (error) {
    return (
      <div style={{ padding: "1rem", background: "rgba(239, 68, 68, 0.1)", border: "1px solid #ef4444", borderRadius: "6px" }}>
        <p style={{ color: "#ef4444", margin: "0 0 0.5rem 0", fontWeight: "bold" }}>Failed to render Mermaid chart</p>
        <pre style={{ margin: 0, overflowX: "auto", fontSize: "0.85rem" }}><code>{chart}</code></pre>
      </div>
    );
  }

  return (
    <div
      className="mermaid-graph"
      style={{
        display: "flex",
        justifyContent: "center",
        padding: "1.5rem",
        background: "var(--surface)",
        border: "1px solid var(--border)",
        borderRadius: "8px",
        overflowX: "auto",
        marginBottom: "1.5rem",
        boxShadow: "0 2px 8px rgba(0,0,0,0.15)",
      }}
      dangerouslySetInnerHTML={{ __html: svg }}
    />
  );
}
