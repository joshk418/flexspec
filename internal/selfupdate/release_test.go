package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

func TestStripLeadingV(t *testing.T) {
	cases := map[string]string{
		"v0.3.5":      "0.3.5",
		"V0.3.5":      "0.3.5",
		"0.3.5":       "0.3.5",
		"":            "",
		"v":           "",
		"v1.2.3-rc.1": "1.2.3-rc.1",
	}
	for in, want := range cases {
		if got := stripLeadingV(in); got != want {
			t.Errorf("stripLeadingV(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestExtractSHA256(t *testing.T) {
	if got := extractSHA256("sha256:abcdef"); got != "abcdef" {
		t.Fatalf("got %q", got)
	}
	if got := extractSHA256("ABCDEF"); got != "abcdef" {
		t.Fatalf("got %q", got)
	}
	if got := extractSHA256(""); got != "" {
		t.Fatalf("got %q", got)
	}
}

func TestParseChecksums(t *testing.T) {
	// sha256sum format: 64-hex-char sum, two spaces, filename.
	body := "1111111111111111111111111111111111111111111111111111111111111111  flexspec_0.3.5_linux_amd64.tar.gz\n" +
		"2222222222222222222222222222222222222222222222222222222222222222  *flexspec_0.3.5_windows_amd64.zip\n" +
		"\n" +
		"garbage line\n" +
		"3333333333333333333333333333333333333333333333333333333333333333  flexspec_0.3.5_darwin_arm64.tar.gz\n"
	m := parseChecksums(body)
	if got, want := m["flexspec_0.3.5_linux_amd64.tar.gz"], "1111111111111111111111111111111111111111111111111111111111111111"; got != want {
		t.Errorf("linux entry = %q, want %q", got, want)
	}
	// binary-mode line has * prefix that should be stripped
	if got, want := m["flexspec_0.3.5_windows_amd64.zip"], "2222222222222222222222222222222222222222222222222222222222222222"; got != want {
		t.Errorf("windows entry = %q, want %q", got, want)
	}
	if _, ok := m["flexspec_0.3.5_darwin_arm64.tar.gz"]; !ok {
		t.Errorf("darwin entry missing: %v", m)
	}
	if _, ok := m["garbage line"]; ok {
		t.Errorf("garbage line should not be in map: %v", m)
	}
	// Short (non-64-char) sums must be rejected.
	if _, ok := m["shortsum"]; ok {
		t.Errorf("short sum should be rejected")
	}
}

func TestReleaseFindAsset(t *testing.T) {
	r := Release{
		Tag:     "v0.3.5",
		Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz"},
			{Name: "flexspec_0.3.5_windows_amd64.zip"},
			{Name: "flexspec_0.3.5_darwin_arm64.tar.gz"},
		},
	}
	got, err := r.FindAsset("windows", "amd64")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "flexspec_0.3.5_windows_amd64.zip" {
		t.Fatalf("got %q", got.Name)
	}

	got, err = r.FindAsset("darwin", "arm64")
	if err != nil {
		t.Fatal(err)
	}
	if got.Name != "flexspec_0.3.5_darwin_arm64.tar.gz" {
		t.Fatalf("got %q", got.Name)
	}

	if _, err := r.FindAsset("freebsd", "amd64"); err == nil {
		t.Fatal("expected error for unknown goos")
	}
}

func TestLatestRelease_parsesTagAndAssets(t *testing.T) {
	orig := httpGet
	t.Cleanup(func() { httpGet = orig })

	payload := releaseJSON{
		TagName: "v0.3.5",
		Assets: []assetJSON{
			{
				Name:               "flexspec_0.3.5_linux_amd64.tar.gz",
				BrowserDownloadURL: "https://example.com/linux.tar.gz",
				Size:               1024,
				Digest:             "sha256:abc",
			},
			{
				Name:               "checksums.txt",
				BrowserDownloadURL: "https://example.com/checksums.txt",
			},
		},
	}
	body, _ := json.Marshal(payload)
	httpGet = func(_ context.Context, _ string) ([]byte, error) {
		return body, nil
	}

	rel, err := LatestRelease(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if rel.Tag != "v0.3.5" || rel.Version != "0.3.5" {
		t.Fatalf("tag=%q version=%q", rel.Tag, rel.Version)
	}
	if len(rel.Assets) != 2 {
		t.Fatalf("assets = %d", len(rel.Assets))
	}
	if rel.Assets[0].DigestSHA256 != "abc" {
		t.Fatalf("digest = %q", rel.Assets[0].DigestSHA256)
	}
}

func TestLatestRelease_emptyTag(t *testing.T) {
	orig := httpGet
	t.Cleanup(func() { httpGet = orig })
	httpGet = func(_ context.Context, _ string) ([]byte, error) {
		return []byte(`{"tag_name":"","assets":[]}`), nil
	}
	if _, err := LatestRelease(context.Background()); err == nil {
		t.Fatal("expected error for empty tag")
	}
}

func TestLatestRelease_httpError(t *testing.T) {
	orig := httpGet
	t.Cleanup(func() { httpGet = orig })
	httpGet = func(_ context.Context, _ string) ([]byte, error) {
		return nil, fmt.Errorf("boom")
	}
	if _, err := LatestRelease(context.Background()); err == nil {
		t.Fatal("expected error")
	}
}

func TestReleaseChecksums(t *testing.T) {
	orig := httpGet
	t.Cleanup(func() { httpGet = orig })

	r := Release{
		Tag: "v0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://example.com/checksums.txt"},
		},
	}
	httpGet = func(_ context.Context, url string) ([]byte, error) {
		if url == "https://example.com/checksums.txt" {
			return []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa  flexspec_0.3.5_linux_amd64.tar.gz\n"), nil
		}
		return nil, fmt.Errorf("unexpected url %s", url)
	}
	m, err := r.Checksums(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if got, want := m["flexspec_0.3.5_linux_amd64.tar.gz"], "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"; got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestReleaseChecksums_missingAsset(t *testing.T) {
	r := Release{Tag: "v0.3.5", Assets: []Asset{{Name: "other.bin"}}}
	if _, err := r.Checksums(context.Background()); err == nil {
		t.Fatal("expected error when checksums.txt asset missing")
	}
}

func TestVersionEqual(t *testing.T) {
	if !versionEqual("0.3.5", "v0.3.5") {
		t.Fatal("0.3.5 == v0.3.5 should be true")
	}
	if versionEqual("0.3.5", "0.3.6") {
		t.Fatal("0.3.5 != 0.3.6 should be true")
	}
}

func TestResumeArgs(t *testing.T) {
	got := ResumeArgs("0.3.4", true, true, true, []string{"a", "b"})
	want := "update --self-update-resume 0.3.4 --skills --migrate --force --only a --only b"
	if strings.Join(got, " ") != want {
		t.Fatalf("got %q, want %q", strings.Join(got, " "), want)
	}

	got = ResumeArgs("0.3.4", false, false, false, nil)
	want = "update --self-update-resume 0.3.4"
	if strings.Join(got, " ") != want {
		t.Fatalf("got %q, want %q", strings.Join(got, " "), want)
	}

	got = ResumeArgs("0.3.4", true, false, false, nil)
	want = "update --self-update-resume 0.3.4 --skills"
	if strings.Join(got, " ") != want {
		t.Fatalf("got %q, want %q", strings.Join(got, " "), want)
	}
}
