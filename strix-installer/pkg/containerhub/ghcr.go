package containerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type GHCRAdapter struct {
	config RegistryEntry
	client *http.Client
}

func NewGHCRAdapter(config RegistryEntry) *GHCRAdapter {
	return &GHCRAdapter{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *GHCRAdapter) Name() string {
	return a.config.Name
}

func (a *GHCRAdapter) Type() RegistryType {
	return RegistryGHCR
}

// FetchImages for GHCR treats the configured URL as a single image repository
// Returns a single "Image" representing the repo, populated with description
func (a *GHCRAdapter) FetchImages(ctx context.Context) ([]Image, error) {
	// GHCR URL format: ghcr.io/owner/image
	parts := strings.Split(strings.TrimPrefix(a.config.URL, "ghcr.io/"), "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid GHCR URL: %s", a.config.URL)
	}

	imageName := parts[1] // e.g. "amd-strix-halo-toolboxes"

	// Start with just the base repo as the "Image"
	// We will fetch tags later on demand
	img := Image{
		Name:        imageName,
		Description: a.config.Description,
		Source:      a.config.Name,
		Registry:    RegistryGHCR,
		URL:         a.config.URL,
	}

	return []Image{img}, nil
}

// GetTags fetches the list of tags for this repository from GHCR API
func (a *GHCRAdapter) GetTags(ctx context.Context, imageName string) ([]Tag, error) {
	// 1. Parse repo path from URL
	// url: ghcr.io/owner/repo
	repoPath := strings.TrimPrefix(a.config.URL, "ghcr.io/")

	// 2. Get auth token
	token, err := a.getToken(ctx, repoPath)
	if err != nil {
		return nil, err
	}

	// 3. List tags
	// https://ghcr.io/v2/owner/repo/tags/list
	apiURL := fmt.Sprintf("https://ghcr.io/v2/%s/tags/list", repoPath)
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list tags: %s", resp.Status)
	}

	var result struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Convert strings to Tag structs
	tags := make([]Tag, len(result.Tags))
	for i, tagName := range result.Tags {
		tags[i] = Tag{
			Name: tagName,
			// For list endpoint, we don't get size/created without extra calls
			// Leaving empty for MVP to improve speed
		}
	}

	return tags, nil
}

func (a *GHCRAdapter) getToken(ctx context.Context, repoPath string) (string, error) {
	// Request anonymous token for public pull
	url := fmt.Sprintf("https://ghcr.io/token?service=ghcr.io&scope=repository:%s:pull", repoPath)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to get auth token: %s", resp.Status)
	}

	var result struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Token, nil
}
