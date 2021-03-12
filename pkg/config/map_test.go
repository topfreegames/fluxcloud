package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMapConfigImplementsConfig(t *testing.T) {
	_ = Config(MapConfig(nil))
}

func TestMapConfigOptional(t *testing.T) {
	m := map[string]string{}
	config := MapConfig(m)

	m["git_url"] = "github.com"
	assert.Equal(t, "github.com", config.Optional("git_url", "hello"))

	delete(m, "git_url")
	assert.Equal(t, "hello", config.Optional("git_url", "hello"))
}

func TestMapConfigRequired(t *testing.T) {
	m := map[string]string{}
	config := MapConfig(m)

	m["git_url"] = "github.com"
	gitURL, err := config.Required("git_url")
	assert.NoError(t, err)
	assert.Equal(t, "github.com", gitURL)

	delete(m, "git_url")
	gitURL, err = config.Required("git_url")
	assert.Error(t, err)
	assert.Equal(t, "", gitURL)
}
