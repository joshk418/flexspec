package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
)

// fakeReleaseServer serves a fixed release JSON plus fake asset bodies.
func fakeReleaseServer(t *testing.T, rel Release, archiveBody []byte, checksumsBody string) {
	t.Helper()
	orig := httpGet
	t.Cleanup(func() { httpGet = orig })

	relJSON, _ := json.Marshal(releaseJSON{
		TagName: rel.Tag,
		Assets:  assetToJSON(rel.Assets),
	})

	httpGet = func(_ context.Context, url string) ([]byte, error) {
		switch url {
		case releasesURL:
			return relJSON, nil
		default:
			for _, a := range rel.Assets {
				if a.DownloadURL == url {
					if strings.HasSuffix(a.Name, ".tar.gz") || strings.HasSuffix(a.Name, ".zip") {
						return archiveBody, nil
					}
					if a.Name == "checksums.txt" {
						return []byte(checksumsBody), nil
					}
				}
			}
			return nil, fmt.Errorf("unexpected url %s", url)
		}
	}
}

func assetToJSON(assets []Asset) []assetJSON {
	out := make([]assetJSON, 0, len(assets))
	for _, a := range assets {
		out = append(out, assetJSON{
			Name:               a.Name,
			BrowserDownloadURL: a.DownloadURL,
			Size:               a.Size,
			Digest:             "sha256:" + a.DigestSHA256,
		})
	}
	return out
}

func TestApplyBinary_alreadyLatest(t *testing.T) {
	fakeReleaseServer(t, Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz"}},
	}, nil, "")

	origSwap := swapBinary
	t.Cleanup(func() { swapBinary = origSwap })
	swapBinary = func(_ []byte) error {
		t.Fatal("swap should not be called when already latest")
		return nil
	}

	res, err := ApplyBinary(context.Background(), "0.3.5", ApplyOpts{})
	if err != nil {
		t.Fatal(err)
	}
	if res.Applied {
		t.Fatal("Applied should be false when already latest")
	}
}

func TestApplyBinary_fullApply(t *testing.T) {
	// Build a real tar.gz containing a "flexspec" binary so ExtractBinary succeeds.
	archiveBody := makeTarGZ(t, "flexspec", "fake new binary")
	// Compute the real sha256 of archiveBody so checksum verification passes.
	wantSum, _ := sha256Hex(archiveBody)
	rel := Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://x/checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz", Size: int64(len(archiveBody))},
		},
	}
	checksumsBody := wantSum + "  flexspec_0.3.5_linux_amd64.tar.gz\n"
	fakeReleaseServer(t, rel, archiveBody, checksumsBody)

	origGoos, origGoarch, origSwap := goos, goarch, swapBinary
	t.Cleanup(func() {
		goos = origGoos
		goarch = origGoarch
		swapBinary = origSwap
	})
	goos, goarch = "linux", "amd64"

	var swappedWith []byte
	swapBinary = func(b []byte) error {
		swappedWith = b
		return nil
	}

	var progress strings.Builder
	res, err := ApplyBinary(context.Background(), "0.3.4", ApplyOpts{Progress: &progress})
	if err != nil {
		t.Fatalf("apply: %v\nprogress: %s", err, progress.String())
	}
	if !res.Applied {
		t.Fatal("Applied should be true")
	}
	if res.ToVersion != "0.3.5" {
		t.Fatalf("ToVersion = %q", res.ToVersion)
	}
	if string(swappedWith) != "fake new binary" {
		t.Fatalf("swappedWith = %q", string(swappedWith))
	}
	if !strings.Contains(progress.String(), "Verifying SHA256... ok") {
		t.Errorf("progress missing verify line:\n%s", progress.String())
	}
	if !strings.Contains(progress.String(), "Executable updated.") {
		t.Errorf("progress missing updated line:\n%s", progress.String())
	}
}

func TestApplyBinary_checksumMismatch(t *testing.T) {
	archiveBody := makeTarGZ(t, "flexspec", "fake binary")
	rel := Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://x/checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz", Size: int64(len(archiveBody))},
		},
	}
	checksumsBody := "0000000000000000000000000000000000000000000000000000000000000000  flexspec_0.3.5_linux_amd64.tar.gz\n"
	fakeReleaseServer(t, rel, archiveBody, checksumsBody)

	origGoos, origGoarch, origSwap := goos, goarch, swapBinary
	t.Cleanup(func() {
		goos = origGoos
		goarch = origGoarch
		swapBinary = origSwap
	})
	goos, goarch = "linux", "amd64"
	swapBinary = func(_ []byte) error {
		t.Fatal("swap should not be called on checksum mismatch")
		return nil
	}

	if _, err := ApplyBinary(context.Background(), "0.3.4", ApplyOpts{}); err == nil {
		t.Fatal("expected checksum mismatch error")
	}
}

func TestApplyBinary_forceReapplies(t *testing.T) {
	archiveBody := makeTarGZ(t, "flexspec", "forced binary")
	wantSum, _ := sha256Hex(archiveBody)
	rel := Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://x/checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz", Size: int64(len(archiveBody))},
		},
	}
	checksumsBody := wantSum + "  flexspec_0.3.5_linux_amd64.tar.gz\n"
	fakeReleaseServer(t, rel, archiveBody, checksumsBody)

	origGoos, origGoarch, origSwap := goos, goarch, swapBinary
	t.Cleanup(func() {
		goos = origGoos
		goarch = origGoarch
		swapBinary = origSwap
	})
	goos, goarch = "linux", "amd64"

	called := false
	swapBinary = func(_ []byte) error {
		called = true
		return nil
	}

	// Current version == latest, but Force=true should still apply.
	res, err := ApplyBinary(context.Background(), "0.3.5", ApplyOpts{Force: true})
	if err != nil {
		t.Fatal(err)
	}
	if !res.Applied || !called {
		t.Fatalf("Applied=%v called=%v, want both true", res.Applied, called)
	}
}

func TestApplyBinary_swapError(t *testing.T) {
	archiveBody := makeTarGZ(t, "flexspec", "fake binary")
	wantSum, _ := sha256Hex(archiveBody)
	rel := Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://x/checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz", Size: int64(len(archiveBody))},
		},
	}
	checksumsBody := wantSum + "  flexspec_0.3.5_linux_amd64.tar.gz\n"
	fakeReleaseServer(t, rel, archiveBody, checksumsBody)

	origGoos, origGoarch, origSwap := goos, goarch, swapBinary
	t.Cleanup(func() {
		goos = origGoos
		goarch = origGoarch
		swapBinary = origSwap
	})
	goos, goarch = "linux", "amd64"
	swapBinary = func(_ []byte) error { return fmt.Errorf("disk full") }

	if _, err := ApplyBinary(context.Background(), "0.3.4", ApplyOpts{}); err == nil {
		t.Fatal("expected swap error")
	}
}

func TestApplyBinary_noMatchingAsset(t *testing.T) {
	rel := Release{
		Tag: "v0.3.5", Version: "0.3.5",
		Assets: []Asset{
			{Name: "checksums.txt", DownloadURL: "https://x/checksums.txt"},
			{Name: "flexspec_0.3.5_linux_amd64.tar.gz", DownloadURL: "https://x/linux.tar.gz"},
		},
	}
	fakeReleaseServer(t, rel, nil, "")

	origGoos, origGoarch, origSwap := goos, goarch, swapBinary
	t.Cleanup(func() {
		goos = origGoos
		goarch = origGoarch
		swapBinary = origSwap
	})
	goos, goarch = "freebsd", "amd64" // not in the asset list
	swapBinary = func(_ []byte) error { t.Fatal("swap should not be called"); return nil }

	if _, err := ApplyBinary(context.Background(), "0.3.4", ApplyOpts{}); err == nil {
		t.Fatal("expected no-asset error")
	}
}
