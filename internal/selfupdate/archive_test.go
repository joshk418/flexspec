package selfupdate

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"
)

// makeTarGZ builds a tar.gz containing a flexspec binary plus a README.
func makeTarGZ(t *testing.T, binaryName, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	files := []struct{ name, body string }{
		{"README.txt", "ignore me"},
		{binaryName, content},
	}
	for _, f := range files {
		hdr := &tar.Header{
			Name: f.name,
			Mode: 0o644,
			Size: int64(len(f.body)),
		}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(f.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

// makeZip builds a zip containing a single binary file (plus a README).
func makeZip(t *testing.T, binaryName, content string) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	files := []struct{ name, body string }{
		{"README.txt", "ignore me"},
		{binaryName, content},
	}
	for _, f := range files {
		w, err := zw.Create(f.name)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := w.Write([]byte(f.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := zw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractBinary_tarGZ(t *testing.T) {
	data := makeTarGZ(t, "flexspec", "binary contents")
	out, err := ExtractBinary(data, "linux", "flexspec_0.3.5_linux_amd64.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "binary contents" {
		t.Fatalf("got %q", string(out))
	}
}

func TestExtractBinary_zip(t *testing.T) {
	data := makeZip(t, "flexspec.exe", "windows binary")
	out, err := ExtractBinary(data, "windows", "flexspec_0.3.5_windows_amd64.zip")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "windows binary" {
		t.Fatalf("got %q", string(out))
	}
}

func TestExtractBinary_unsupportedFormat(t *testing.T) {
	if _, err := ExtractBinary([]byte("x"), "linux", "flexspec.rar"); err == nil {
		t.Fatal("expected error for unsupported format")
	}
}

func TestExtractBinary_missingBinary(t *testing.T) {
	data := makeTarGZ(t, "wrongname", "content")
	if _, err := ExtractBinary(data, "linux", "flexspec_0.3.5_linux_amd64.tar.gz"); err == nil {
		t.Fatal("expected error when binary not found in archive")
	}
}

func TestExtractBinary_pathWithSubdir(t *testing.T) {
	// Be lenient if the binary is nested instead of at the archive root.
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	files := []struct{ name, body string }{
		{"flexspec-v0.3.5/flexspec", "nested"},
	}
	for _, f := range files {
		hdr := &tar.Header{Name: f.name, Mode: 0o644, Size: int64(len(f.body))}
		if err := tw.WriteHeader(hdr); err != nil {
			t.Fatal(err)
		}
		if _, err := tw.Write([]byte(f.body)); err != nil {
			t.Fatal(err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gw.Close(); err != nil {
		t.Fatal(err)
	}
	out, err := ExtractBinary(buf.Bytes(), "linux", "flexspec_0.3.5_linux_amd64.tar.gz")
	if err != nil {
		t.Fatal(err)
	}
	if string(out) != "nested" {
		t.Fatalf("got %q", string(out))
	}
}

func TestNameBase(t *testing.T) {
	cases := map[string]string{
		"flexspec":           "flexspec",
		"dir/flexspec":       "flexspec",
		"dir\\flexspec":      "flexspec",
		"a/b/c/flexspec.exe": "flexspec.exe",
	}
	for in, want := range cases {
		if got := nameBase(in); got != want {
			t.Errorf("nameBase(%q) = %q, want %q", in, got, want)
		}
	}
}
