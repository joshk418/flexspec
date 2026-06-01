package ui

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func writeUIProject(t *testing.T, root string) {
	t.Helper()
	base := filepath.Join(root, ".flexspec")
	if err := os.MkdirAll(base, 0o755); err != nil {
		t.Fatal(err)
	}
	cfg := "specs_dir: specs\nalways_one_shot: false\nspec_template:\n"
	if err := os.WriteFile(filepath.Join(base, "config.yaml"), []byte(cfg), 0o644); err != nil {
		t.Fatal(err)
	}
	specDir := filepath.Join(root, "specs", "001-test")
	if err := os.MkdirAll(specDir, 0o755); err != nil {
		t.Fatal(err)
	}
	readme := "---\nname: Test\ndescription: d\nstatus: planned\nspec_type: simple\n---\n\n# Body\n"
	if err := os.WriteFile(filepath.Join(specDir, "README.md"), []byte(readme), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestServer_healthAndSpecs(t *testing.T) {
	root := t.TempDir()
	writeUIProject(t, root)

	srv, err := NewServer(root, "127.0.0.1", 0, StubStaticFS())
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(srv.http.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/health")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("health status = %d", resp.StatusCode)
	}

	resp2, err := http.Get(ts.URL + "/api/specs")
	if err != nil {
		t.Fatal(err)
	}
	defer resp2.Body.Close()
	var specs []SpecJSON
	if err := json.NewDecoder(resp2.Body).Decode(&specs); err != nil {
		t.Fatal(err)
	}
	if len(specs) != 1 || specs[0].Name != "Test" {
		t.Fatalf("specs = %+v", specs)
	}
}

func TestNewServer_missingConfig(t *testing.T) {
	root := t.TempDir()
	_, err := NewServer(root, "127.0.0.1", 3000, StubStaticFS())
	if err == nil {
		t.Fatal("expected error for missing config")
	}
}

func TestServer_configGetPut(t *testing.T) {
	root := t.TempDir()
	writeUIProject(t, root)

	srv, err := NewServer(root, "127.0.0.1", 0, StubStaticFS())
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(srv.http.Handler)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/config")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET config status = %d", resp.StatusCode)
	}

	var got map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got["specs_dir"] != "specs" || got["always_one_shot"] != false {
		t.Fatalf("config = %+v", got)
	}

	body := []byte(`{"specs_dir":"docs-specs","always_one_shot":true,"spec_template":"simple"}`)
	putResp, err := http.NewRequest(http.MethodPut, ts.URL+"/api/config", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	putResp.Header.Set("Content-Type", "application/json")
	putResult, err := http.DefaultClient.Do(putResp)
	if err != nil {
		t.Fatal(err)
	}
	defer putResult.Body.Close()
	if putResult.StatusCode != http.StatusOK {
		t.Fatalf("PUT config status = %d", putResult.StatusCode)
	}

	var saved map[string]any
	if err := json.NewDecoder(putResult.Body).Decode(&saved); err != nil {
		t.Fatal(err)
	}
	if saved["specs_dir"] != "docs-specs" || saved["always_one_shot"] != true || saved["spec_template"] != "simple" {
		t.Fatalf("saved = %+v", saved)
	}

	badBody := []byte(`{"specs_dir":"","always_one_shot":false,"spec_template":""}`)
	badReq, err := http.NewRequest(http.MethodPut, ts.URL+"/api/config", bytes.NewReader(badBody))
	if err != nil {
		t.Fatal(err)
	}
	badReq.Header.Set("Content-Type", "application/json")
	badResp, err := http.DefaultClient.Do(badReq)
	if err != nil {
		t.Fatal(err)
	}
	defer badResp.Body.Close()
	if badResp.StatusCode != http.StatusBadRequest {
		t.Fatalf("invalid PUT status = %d, want 400", badResp.StatusCode)
	}
}
