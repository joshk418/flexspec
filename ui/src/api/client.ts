export type Spec = {
  id: string;
  dir: string;
  name: string;
  description: string;
  status: string;
  spec_type: string;
  tasks?: Task[];
};

export type Task = {
  id: string;
  file: string;
  name: string;
  status: string;
};

export type SpecDetail = Omit<Spec, "tasks"> & {
  markdown: string;
  tasks: TaskDetail[];
};

export type TaskDetail = Task & {
  markdown: string;
};

export type ProjectConfig = {
  specs_dir: string;
  always_one_shot: boolean;
  spec_template: string;
};

const base = "";

export async function fetchSpecs(): Promise<Spec[]> {
  const res = await fetch(`${base}/api/specs`);
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function fetchSpecDetail(dir: string): Promise<SpecDetail> {
  const res = await fetch(`${base}/api/specs/${encodeURIComponent(dir)}`);
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function fetchConfig(): Promise<ProjectConfig> {
  const res = await fetch(`${base}/api/config`);
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

export async function saveConfig(config: ProjectConfig): Promise<ProjectConfig> {
  const res = await fetch(`${base}/api/config`, {
    method: "PUT",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(config),
  });
  if (!res.ok) throw new Error(await errorMessage(res));
  return res.json();
}

async function errorMessage(res: Response): Promise<string> {
  try {
    const body = await res.json();
    if (body?.error) return body.error;
  } catch {
    /* ignore */
  }
  return res.statusText;
}

export function subscribeEvents(onChange: () => void): () => void {
  const es = new EventSource(`${base}/api/events`);
  es.addEventListener("specs-changed", () => onChange());
  return () => es.close();
}

export const SPEC_COLUMNS = [
  "initial",
  "refined",
  "planned",
  "in_progress",
  "in_review",
  "complete",
] as const;

export function columnForStatus(status: string): string {
  const s = status.trim().toLowerCase();
  if ((SPEC_COLUMNS as readonly string[]).includes(s)) return s;
  return "unassigned";
}
