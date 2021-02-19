package config

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChainConfigImplementsConfig(t *testing.T) {
	_ = Config(&ChainConfig{})
}

func TestChainConfigOptional(t *testing.T) {
	for name, test := range map[string]struct {
		head, tail                Config
		key, defaultValue, output string
	}{
		"present-in-head": {
			head: MapConfig(map[string]string{
				"present-key": "present-value",
			}),
			tail:         MapConfig(map[string]string{}),
			key:          "present-key",
			defaultValue: "wrong-value",
			output:       "present-value",
		},
		"present-in-tail": {
			head: MapConfig(map[string]string{}),
			tail: MapConfig(map[string]string{
				"present-key": "present-value",
			}),
			key:          "present-key",
			defaultValue: "wrong-value",
			output:       "present-value",
		},
		"not-present": {
			head:         MapConfig(map[string]string{}),
			tail:         MapConfig(map[string]string{}),
			key:          "not-present-key",
			defaultValue: "default-value",
			output:       "default-value",
		},
	} {
		t.Run(name, func(t *testing.T) {
			chain := NewChain(test.head, test.tail)
			output := chain.Optional(test.key, test.defaultValue)
			assert.Equal(t, test.output, output)
		})
	}
}

func TestChainConfigRequired(t *testing.T) {
	for name, test := range map[string]struct {
		head, tail  Config
		key, output string
		err         error
	}{
		"present-in-head": {
			head: MapConfig(map[string]string{
				"present-key": "present-value",
			}),
			tail:   MapConfig(map[string]string{}),
			key:    "present-key",
			output: "present-value",
		},
		"present-in-tail": {
			head: MapConfig(map[string]string{}),
			tail: MapConfig(map[string]string{
				"present-key": "present-value",
			}),
			key:    "present-key",
			output: "present-value",
		},
		"not-present": {
			head: MapConfig(map[string]string{}),
			tail: MapConfig(map[string]string{}),
			key:  "not-present-key",
			err:  errors.New("key not found: not-present-key"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			chain := NewChain(test.head, test.tail)
			output, err := chain.Required(test.key)
			assert.Equal(t, test.err, err)
			assert.Equal(t, test.output, output)
		})
	}
}
