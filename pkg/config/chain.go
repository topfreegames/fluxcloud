package config

// ChainConfig fetches configurations keys from a head Config, and fallbacks
// to a tail Config if the configuration is not present.
type ChainConfig struct {
	head, tail Config
}

// NewChain creates a ChainConfig given a list of Configs.
// The function panics if an empty list is passed.
func NewChain(configs ...Config) Config {
	if len(configs) == 0 {
		panic("cannot create config chain from empty list")
	}

	if len(configs) == 1 {
		return configs[0]
	}

	return &ChainConfig{
		head: configs[0],
		tail: NewChain(configs[1:]...),
	}
}

func (c *ChainConfig) Optional(key string, defaultValue string) string {
	if s, err := c.head.Required(key); err == nil {
		return s
	}
	return c.tail.Optional(key, defaultValue)
}

func (c *ChainConfig) Required(key string) (string, error) {
	if s, err := c.head.Required(key); err == nil {
		return s, nil
	}
	return c.tail.Required(key)
}
