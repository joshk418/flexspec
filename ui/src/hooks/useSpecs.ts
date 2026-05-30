import { useCallback, useEffect, useState } from "react";
import { fetchSpecs, subscribeEvents, type Spec } from "../api/client";

export function useSpecs() {
  const [specs, setSpecs] = useState<Spec[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const reload = useCallback(async () => {
    try {
      setError(null);
      const data = await fetchSpecs();
      setSpecs(data);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Failed to load specs");
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    reload();
    return subscribeEvents(reload);
  }, [reload]);

  return { specs, loading, error, reload };
}
