package selfupdate

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// githubOwner and githubRepo identify the GitHub release source.
const (
	githubOwner = "joshk418"
	githubRepo  = "flexspec"
	releasesURL = "https://api.github.com/repos/" + githubOwner + "/" + githubRepo + "/releases/latest"
)

// Asset is one downloadable file attached to a GitHub release.
type Asset struct {
	Name         string
	DownloadURL  string
	Size         int64
	DigestSHA256 string // hex sha256 from the GitHub asset "digest" field, if present
	BrowserURL   string
}

// Release describes the latest GitHub release for the flexspec repo.
type Release struct {
	Tag     string // e.g. "v0.3.5"
	Version string // e.g. "0.3.5" (Tag without leading 'v')
	Assets  []Asset
}

// assetJSON mirrors the subset of the GitHub releases API asset shape we use.
type assetJSON struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	Digest             string `json:"digest"` // "sha256:hex..."
}

type releaseJSON struct {
	TagName string      `json:"tag_name"`
	Assets  []assetJSON `json:"assets"`
}

// httpGet is the HTTP client hook; tests inject a fake transport.
var httpGet func(ctx context.Context, url string) ([]byte, error) = defaultHTTPGet

func defaultHTTPGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return nil, fmt.Errorf("GET %s: %s: %s", url, resp.Status, strings.TrimSpace(string(body)))
	}
	return io.ReadAll(resp.Body)
}

// LatestRelease fetches the latest GitHub release metadata (no asset bodies).
func LatestRelease(ctx context.Context) (Release, error) {
	body, err := httpGet(ctx, releasesURL)
	if err != nil {
		return Release{}, err
	}
	var rj releaseJSON
	if err := json.Unmarshal(body, &rj); err != nil {
		return Release{}, fmt.Errorf("parse release JSON: %w", err)
	}
	if rj.TagName == "" {
		return Release{}, fmt.Errorf("release has empty tag_name")
	}
	r := Release{
		Tag:     rj.TagName,
		Version: stripLeadingV(rj.TagName),
	}
	for _, a := range rj.Assets {
		asset := Asset{
			Name:         a.Name,
			DownloadURL:  a.BrowserDownloadURL,
			Size:         a.Size,
			BrowserURL:   a.BrowserDownloadURL,
			DigestSHA256: extractSHA256(a.Digest),
		}
		r.Assets = append(r.Assets, asset)
	}
	return r, nil
}

// FindAsset returns the asset matching goos/goarch (e.g. "windows","amd64" -> "flexspec_0.3.5_windows_amd64.zip").
func (r Release) FindAsset(goos, goarch string) (Asset, error) {
	want := fmt.Sprintf("_%s_%s.", goos, goarch)
	for _, a := range r.Assets {
		if strings.Contains(a.Name, want) {
			return a, nil
		}
	}
	return Asset{}, fmt.Errorf("no asset for %s/%s in release %s (have: %s)", goos, goarch, r.Tag, r.assetNames())
}

func (r Release) assetNames() string {
	names := make([]string, 0, len(r.Assets))
	for _, a := range r.Assets {
		names = append(names, a.Name)
	}
	return strings.Join(names, ", ")
}

// Checksums downloads checksums.txt and returns a map of filename -> lowercase hex sha256.
func (r Release) Checksums(ctx context.Context) (map[string]string, error) {
	url := ""
	for _, a := range r.Assets {
		if a.Name == "checksums.txt" {
			url = a.DownloadURL
			break
		}
	}
	if url == "" {
		return nil, fmt.Errorf("release %s has no checksums.txt asset", r.Tag)
	}
	body, err := httpGet(ctx, url)
	if err != nil {
		return nil, err
	}
	return parseChecksums(string(body)), nil
}

// parseChecksums parses "<hex>  <filename>" lines (sha256sum format); non-matching lines are ignored.
func parseChecksums(s string) map[string]string {
	out := make(map[string]string)
	for _, line := range strings.Split(s, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Split on first whitespace; trim filename.
		idx := strings.IndexAny(line, " \t")
		if idx <= 0 {
			continue
		}
		sum := strings.TrimSpace(line[:idx])
		name := strings.TrimSpace(line[idx:])
		// sha256sum prefixes binary-mode lines with "*"; strip it.
		name = strings.TrimPrefix(name, "*")
		if len(sum) == 64 {
			out[name] = strings.ToLower(sum)
		}
	}
	return out
}

// stripLeadingV returns tag without a leading 'v' or 'V'.
func stripLeadingV(tag string) string {
	if len(tag) > 0 && (tag[0] == 'v' || tag[0] == 'V') {
		return tag[1:]
	}
	return tag
}

// extractSHA256 parses a "sha256:hex" digest string and returns the hex part.
func extractSHA256(digest string) string {
	if digest == "" {
		return ""
	}
	if i := strings.Index(digest, ":"); i >= 0 {
		return strings.ToLower(strings.TrimSpace(digest[i+1:]))
	}
	return strings.ToLower(strings.TrimSpace(digest))
}
