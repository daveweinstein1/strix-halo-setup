package containerhub

import (
	"context"
	"time"
)

// RegistryType defines the source of the registry
type RegistryType string

const (
	RegistryGHCR RegistryType = "ghcr"
	RegistryJSON RegistryType = "json"
)

// Image represents a single installable container image
type Image struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Source      string       `json:"source"` // "kyuz0", "amd", etc.
	Tags        []Tag        `json:"tags"`
	Registry    RegistryType `json:"registry_type"`
	URL         string       `json:"url"` // Full pull URL
}

// Tag represents a specific version of an image
type Tag struct {
	Name      string    `json:"name"`
	Digest    string    `json:"digest"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// Registry defines the interface for fetching images from a source
type Registry interface {
	Name() string
	Type() RegistryType
	FetchImages(ctx context.Context) ([]Image, error)
	GetTags(ctx context.Context, imageName string) ([]Tag, error)
}

// Config represents the registries.yaml configuration
type RegistryConfig struct {
	Registries []RegistryEntry `yaml:"registries"`
}

type RegistryEntry struct {
	Name        string       `yaml:"name"`
	Type        RegistryType `yaml:"type"`
	URL         string       `yaml:"url"`
	Description string       `yaml:"description"`
	Priority    int          `yaml:"priority"`
}
