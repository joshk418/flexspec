package selfupdate

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"runtime"

	selfupdatelib "github.com/minio/selfupdate"
)

// ApplyOpts controls ApplyBinary behavior.
type ApplyOpts struct {
	// Force re-downloads and re-applies even when currentVersion == latest.
	Force bool
	// Progress receives status lines during the update. May be nil.
	Progress io.Writer
}

// ApplyResult describes what happened during an ApplyBinary call.
type ApplyResult struct {
	Release     Release
	Asset       Asset
	FromVersion string
	ToVersion   string
	Applied     bool // true if the binary was actually swapped
}

// ApplyBinary fetches the latest release, then downloads, verifies, and swaps the binary when needed.
func ApplyBinary(ctx context.Context, currentVersion string, opts ApplyOpts) (ApplyResult, error) {
	progress := opts.Progress
	res := ApplyResult{FromVersion: currentVersion}

	rel, err := LatestRelease(ctx)
	if err != nil {
		return res, err
	}
	res.Release = rel

	if !opts.Force && versionEqual(currentVersion, rel.Version) {
		writeProgress(progress, "Already on latest (v%s).\n", rel.Version)
		return res, nil
	}

	asset, err := rel.FindAsset(goos, goarch)
	if err != nil {
		return res, err
	}
	res.Asset = asset
	res.ToVersion = rel.Version

	writeProgress(progress, "Latest: v%s (current v%s)\n", rel.Version, currentVersion)
	writeProgress(progress, "Downloading %s (%.1f MB)...\n", asset.Name, float64(asset.Size)/1024/1024)

	archiveBytes, err := downloadAsset(ctx, asset)
	if err != nil {
		return res, err
	}

	// Verify SHA256 against checksums.txt before touching the running binary.
	checksums, err := rel.Checksums(ctx)
	if err != nil {
		return res, fmt.Errorf("verify checksum: %w", err)
	}
	want, ok := checksums[asset.Name]
	if !ok {
		// Fall back to the asset's embedded digest if checksums.txt is missing the entry.
		if asset.DigestSHA256 == "" {
			return res, fmt.Errorf("no sha256 for %s in checksums.txt and no asset digest", asset.Name)
		}
		want = asset.DigestSHA256
	}
	got := sha256Hex(archiveBytes)
	if got != want {
		return res, fmt.Errorf("checksum mismatch for %s: want %s, got %s", asset.Name, want, got)
	}
	writeProgress(progress, "Verifying SHA256... ok\n")

	binary, err := ExtractBinary(archiveBytes, goos, asset.Name)
	if err != nil {
		return res, fmt.Errorf("extract binary: %w", err)
	}

	writeProgress(progress, "Replacing executable...\n")
	// swapBinary swaps the running executable atomically (Windows-safe: renames in-use binary to *.old).
	if err := swapBinary(binary); err != nil {
		return res, err
	}

	res.Applied = true
	writeProgress(progress, "Executable updated.\n")
	return res, nil
}

// downloadAsset fetches the full asset body.
func downloadAsset(ctx context.Context, asset Asset) ([]byte, error) {
	body, err := httpGet(ctx, asset.DownloadURL)
	if err != nil {
		return nil, fmt.Errorf("download %s: %w", asset.Name, err)
	}
	return body, nil
}

// sha256Hex returns the lowercase hex sha256 of data.
func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// versionEqual compares two version strings, tolerating a leading 'v'.
func versionEqual(a, b string) bool {
	return stripLeadingV(a) == stripLeadingV(b)
}

// writeProgress writes a formatted line to w if non-nil; errors are ignored (best-effort).
func writeProgress(w io.Writer, format string, args ...any) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintf(w, format, args...)
}

// swapBinary is the seam over selfupdatelib.Apply so tests don't replace the test binary.
var swapBinary = swapBinaryDefault

// swapBinaryDefault calls minio/selfupdate to atomically swap the running executable.
func swapBinaryDefault(binary []byte) error {
	applyOpts := selfupdatelib.Options{
		TargetMode: 0o755,
	}
	if err := selfupdatelib.Apply(bytes.NewReader(binary), applyOpts); err != nil {
		if rerr := selfupdatelib.RollbackError(err); rerr != nil {
			return fmt.Errorf("apply update (rollback also failed: %v): %w", rerr, err)
		}
		return fmt.Errorf("apply update: %w", err)
	}
	return nil
}

// goos and goarch are seams over runtime.GOOS/GOARCH for testing.
var (
	goos   = runtime.GOOS
	goarch = runtime.GOARCH
)
