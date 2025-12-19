package helper

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	// Default TTL untuk cache (dalam detik)
	DefaultCacheTTL = 30 * time.Minute
	// TTL untuk pohon kinerja (lebih lama karena data jarang berubah)
	PohonKinerjaCacheTTL = 5 * time.Minute
)

// CacheKeyPrefix untuk berbagai jenis cache
const (
	CacheKeyPohonKinerjaOpdAll = "pohon_kinerja_opd_all"
)

// GenerateCacheKey menghasilkan cache key berdasarkan prefix dan parameter
func GenerateCacheKey(prefix string, params ...string) string {
	key := prefix
	for _, param := range params {
		key = fmt.Sprintf("%s:%s", key, param)
	}
	return key
}

// GetFromCache mengambil data dari Redis cache
func GetFromCache(ctx context.Context, client *redis.Client, key string, dest interface{}) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return fmt.Errorf("cache miss")
		}
		return fmt.Errorf("error getting from cache: %w", err)
	}

	err = json.Unmarshal([]byte(val), dest)
	if err != nil {
		return fmt.Errorf("error unmarshaling cache: %w", err)
	}

	return nil
}

// SetToCache menyimpan data ke Redis cache
func SetToCache(ctx context.Context, client *redis.Client, key string, value interface{}, ttl time.Duration) error {
	if client == nil {
		return fmt.Errorf("redis client is nil")
	}

	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling to cache: %w", err)
	}

	err = client.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("error setting cache: %w", err)
	}

	return nil
}

// DeleteCache menghapus data dari cache
func DeleteCache(ctx context.Context, client *redis.Client, key string) error {
	if client == nil {
		return nil // Ignore if Redis is not available
	}

	err := client.Del(ctx, key).Err()
	if err != nil {
		return fmt.Errorf("error deleting cache: %w", err)
	}

	return nil
}

// DeleteCacheByPattern menghapus semua cache yang match pattern (misal: pohon_kinerja_opd_all:*)
func DeleteCacheByPattern(ctx context.Context, client *redis.Client, pattern string) error {
	if client == nil {
		return nil // Ignore if Redis is not available
	}

	iter := client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		err := client.Del(ctx, iter.Val()).Err()
		if err != nil {
			return fmt.Errorf("error deleting cache key %s: %w", iter.Val(), err)
		}
	}

	if err := iter.Err(); err != nil {
		return fmt.Errorf("error scanning cache: %w", err)
	}

	return nil
}
