package persist

import "github.com/pkg/errors"

type Cache interface {
	// Get returns the value stored under the given key, if any
	// returns zero-value (empty string) if no value found
	Get(key string) string

	// Put stores the given value under the given key
	// empty key or value yields and error
	Put(key string, value string) error
}

type mockCache struct {
	// messages is a map of digests to plaintext messages
	values map[string]string
}

func (c *mockCache) Get(key string) string {
	message, ok := c.values[key]
	if !ok {
		// Cache miss, return zero-value
		return ""
	}

	return message
}

func (c *mockCache) Put(key string, value string) error {
	if len(key) < 1 || len(value) < 1 {
		return errors.Errorf("can't put {%s: %s}, non-empty key and value required", key, value)
	}

	c.values[key] = value

	return nil
}

func NewMockCache() *mockCache {
	return &mockCache{
		map[string]string{},
	}
}
