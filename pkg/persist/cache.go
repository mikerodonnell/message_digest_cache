package persist

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
)

// see docker-compose.yml
// set this to :6379 if redis is running on localhost
const redisHost = "redis:6379"

type Cache interface {
	// Get returns the value stored under the given key, if any
	// returns zero-value (empty string) if no value found
	Get(key string) string

	// Put stores the given value under the given key
	// empty key or value yields and error
	Put(key string, value string) error
}

type localCache struct {
	// messages is a map of digests to plaintext messages
	values map[string]string
}

func (c *localCache) Get(key string) string {
	message, ok := c.values[key]
	if !ok {
		// cache miss, return zero-value
		return ""
	}

	return message
}

func (c *localCache) Put(key string, value string) error {
	if len(key) < 1 || len(value) < 1 {
		return errors.Errorf("can't put {%s: %s}, non-empty key and value required", key, value)
	}

	c.values[key] = value

	return nil
}

func NewLocalCache() *localCache {
	return &localCache{
		map[string]string{},
	}
}

type redisCache struct {
	redis.Conn
}

func (c *redisCache) Get(key string) string {
	log.Println(fmt.Sprintf("fetching %s from redis", key))
	response, err := c.Do("GET", key)
	if err != nil {
		// this should never happen; if the key doesn't exist redis returns an empty response, no error
		// log and swallow this error for this prototype app
		log.Println(fmt.Sprintf("unexpected error getting %s, returning zero value", key))
		return ""
	}

	value, ok := response.([]byte)
	if value == nil || !ok {
		return ""
	}

	return string(value)
}

func (c *redisCache) Put(key string, value string) error {
	if len(key) < 1 || len(value) < 1 {
		return errors.Errorf("can't put {%s: %s}, non-empty key and value required", key, value)
	}

	_, err := c.Do("SET", key, value)
	if err != nil {
		return errors.Wrapf(err, "setting value %s for key %s", value, key)
	}

	return nil
}

// NewRedisCache instantiates a new redisCache with an open connection
// caller is responsible for calling Close() when finished
func NewRedisCache() (*redisCache, error) {
	connection, err := redis.Dial("tcp", redisHost)
	if err != nil {
		return nil, errors.Wrap(err, "instantiating new redis cache")
	}

	return &redisCache{connection}, nil
}
