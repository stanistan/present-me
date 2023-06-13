package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/peterbourgon/diskv/v3"

	"github.com/stanistan/present-me/internal/log"
)

type Cache struct {
	disabled bool
	d        *diskv.Diskv
}

type Options struct {
	TTL          time.Duration
	ForceRefresh bool
}

type DataKey struct {
	Prefix  string
	Hashing any
}

func DataKeyFor(prefix string, hashing ...any) DataKey {
	return DataKey{
		Prefix:  prefix,
		Hashing: hashing,
	}
}

func (k *DataKey) String() string {
	return fmt.Sprintf("%s-%v", k.Prefix, k.Hashing)
}

type Provider[T any] struct {
	DataKey
	Fetch func() (T, error)
}

func Apply[T any](ctx context.Context, k DataKey, f func() (T, error)) (T, error) {
	var (
		out          T
		ttl          time.Duration
		forceRefresh bool
	)

	c := Ctx(ctx)
	p := Provider[T]{DataKey: k, Fetch: f}
	opts, ok := OptionsFromContext(ctx)
	if ok {
		ttl = opts.TTL
		forceRefresh = opts.ForceRefresh
	}

	if !forceRefresh {
		ok, err := c.Read(ctx, p.DataKey, &out, ttl)
		if err != nil {
			return out, err
		}
		if ok {
			return out, nil
		}
	}

	data, err := p.Fetch()
	if err != nil {
		return out, err
	}

	err = c.Write(p.DataKey, data)
	if err != nil {
		return out, err
	}

	return data, nil
}

func (c *Cache) Read(ctx context.Context, key DataKey, into any, ttl time.Duration) (bool, error) {
	if c.disabled {
		return false, nil
	}

	k, err := Key(key)
	if err != nil {
		return false, err
	}

	bytes, err := c.d.Read(k)
	if err != nil {
		log.Ctx(ctx).Warn().Err(err).Msg("")
		return false, nil // FILE MISSING, Do a check here:)
	}

	storedAt, err := Unmarshal(bytes, into)
	if err != nil {
		return false, err
	}

	if storedAt == nil || time.Since(*storedAt) > ttl {
		log.Ctx(ctx).Info().Str("key", key.String()).Msg("data expired")
		return false, nil
	}

	return true, nil
}

func (c *Cache) Write(key DataKey, data any) error {
	if c.disabled {
		return nil
	}

	k, err := Key(key)
	if err != nil {
		return err
	}

	bytes, err := Marshal(data)
	if err != nil {
		return err
	}

	return c.d.Write(k, bytes)
}

type CacheOpts struct {
	Enabled        bool   `name:"enabled" env:"DISK_CACHE_ENABLED"`
	BasePath       string `name:"base-path" env:"DISK_CACHE_BASE_PATH"`
	CacheMaxSizeKB uint64 `name:"cache-max-size" env:"DISK_CACHE_MAX_SIZE_KB" default:"1024"`
}

func NewCache(ctx context.Context, opts CacheOpts) *Cache {
	if !opts.Enabled {
		return &Cache{disabled: true}
	}

	cacheOpts := diskv.Options{
		BasePath:     opts.BasePath,
		CacheSizeMax: opts.CacheMaxSizeKB * 1024,
	}

	log.Ctx(ctx).Info().Msgf("initializing cache basePath=%s size=%d", cacheOpts.BasePath, cacheOpts.CacheSizeMax)
	return &Cache{d: diskv.New(cacheOpts)}
}

type Value struct {
	At   time.Time
	Data json.RawMessage
}

func Key(k DataKey) (string, error) {
	hash, err := hashstructure.Hash(k.Hashing, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s.json", k.Prefix, strconv.FormatUint(hash, 10)), nil
}

func Marshal(data any) ([]byte, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(Value{At: time.Now(), Data: dataBytes})
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func Unmarshal(bytes []byte, into any) (*time.Time, error) {
	var value Value
	if err := json.Unmarshal(bytes, &value); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(value.Data, into); err != nil {
		return nil, err
	}

	return &value.At, nil
}
