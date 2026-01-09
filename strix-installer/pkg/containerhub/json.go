package containerhub

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type JSONAdapter struct {
	config RegistryEntry
	client *http.Client
}

func NewJSONAdapter(config RegistryEntry) *JSONAdapter {
	return &JSONAdapter{
		config: config,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (a *JSONAdapter) Name() string {
	return a.config.Name
}

func (a *JSONAdapter) Type() RegistryType {
	return RegistryJSON
}

type jsonRegistryFormat struct {
	Images []jsonImage `json:"images"`
}

type jsonImage struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	URL         string   `json:"url"`
	Tags        []string `json:"tags"`
}

func (a *JSONAdapter) FetchImages(ctx context.Context) ([]Image, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", a.config.URL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := a.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch registry JSON: %s", resp.Status)
	}

	var data jsonRegistryFormat
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("invalid registry JSON: %v", err)
	}

	images := make([]Image, len(data.Images))
	for i, jImg := range data.Images {
		tags := make([]Tag, len(jImg.Tags))
		for k, t := range jImg.Tags {
			tags[k] = Tag{Name: t}
		}

		images[i] = Image{
			Name:        jImg.Name,
			Description: jImg.Description,
			Source:      a.config.Name,
			Registry:    RegistryJSON,
			URL:         jImg.URL,
			Tags:        tags,
		}
	}

	return images, nil
}

func (a *JSONAdapter) GetTags(ctx context.Context, imageName string) ([]Tag, error) {
	// For JSON registry, tags are usually pre-fetched in FetchImages
	// But if we needed to refetch, we'd pull the whole JSON again.
	// For efficiency, we rely on the cached image data from FetchImages ideally,
	// but the interface requires this method.

	// Re-fetch logic (simple but inefficient)
	images, err := a.FetchImages(ctx)
	if err != nil {
		return nil, err
	}

	for _, img := range images {
		if img.Name == imageName {
			return img.Tags, nil
		}
	}

	return nil, fmt.Errorf("image not found: %s", imageName)
}
