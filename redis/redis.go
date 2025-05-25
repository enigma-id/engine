// redis.go
package redis

import (
	"encoding/json"

	"github.com/gomodule/redigo/redis"
)

// Redis wraps a Redis connection pool and key prefix.
type Redis struct {
	Prefix string      // Prefix for all Redis keys
	Pool   *redis.Pool // Connection pool
}

// Save stores the value as JSON under the given key.
func (r *Redis) Save(key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	conn := r.Pool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", r.key(key), data)
	return err
}

// Read retrieves and unmarshals the value stored under the given key.
func (r *Redis) Read(key string, out any) error {
	conn := r.Pool.Get()
	defer conn.Close()

	data, err := redis.Bytes(conn.Do("GET", r.key(key)))
	if err != nil {
		return err
	}

	return json.Unmarshal(data, out)
}

// Delete removes the specified key from Redis.
func (r *Redis) Delete(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", r.key(key))
	return err
}

// ping checks if the Redis connection is alive.
func (r *Redis) ping() error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := redis.String(conn.Do("PING"))
	return err
}

// key applies the configured prefix to the given key.
func (r *Redis) key(k string) string {
	return r.Prefix + k
}
