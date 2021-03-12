package config

import "fmt"

// MapConfig fetches configuration keys from a map.
type MapConfig map[string]string

func (m MapConfig) Optional(key string, defaultValue string) string {
	if v, ok := m[key]; ok {
		return v
	}
	return defaultValue
}

func (m MapConfig) Required(key string) (string, error) {
	if v, ok := m[key]; ok {
		return v, nil
	}
	return "", fmt.Errorf("key not found: %s", key)
}
