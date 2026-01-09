package containerhub

import (
	"context"
	"fmt"
	"os"
	"sort"
	"sync"

	"gopkg.in/yaml.v3"
)

type Manager struct {
	registries []Registry
}

func NewManager() *Manager {
	return &Manager{
		registries: make([]Registry, 0),
	}
}

// LoadConfigFromPath loads registries from a YAML file
func (m *Manager) LoadConfigFromPath(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config RegistryConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return err
	}

	for _, entry := range config.Registries {
		switch entry.Type {
		case RegistryGHCR:
			m.registries = append(m.registries, NewGHCRAdapter(entry))
		case RegistryJSON:
			m.registries = append(m.registries, NewJSONAdapter(entry))
		default:
			// Log unknown type? Return error?
			// For now, skip
		}
	}

	return nil
}

// FetchAllImages queries all configured registries concurrently
func (m *Manager) FetchAllImages(ctx context.Context) ([]Image, error) {
	var wg sync.WaitGroup
	results := make(chan []Image, len(m.registries))
	errors := make(chan error, len(m.registries))

	for _, reg := range m.registries {
		wg.Add(1)
		go func(r Registry) {
			defer wg.Done()
			imgs, err := r.FetchImages(ctx)
			if err != nil {
				errors <- fmt.Errorf("registry %s failed: %v", r.Name(), err)
				return
			}
			results <- imgs
		}(reg)
	}

	wg.Wait()
	close(results)
	close(errors)

	// Collect results
	var allImages []Image

	// Check errors (logging them would be better than returning 1 failure)
	// For MVP, if one fails, we still return others

	for imgs := range results {
		allImages = append(allImages, imgs...)
	}

	// Sort by Name
	sort.Slice(allImages, func(i, j int) bool {
		return allImages[i].Name < allImages[j].Name
	})

	return allImages, nil
}
