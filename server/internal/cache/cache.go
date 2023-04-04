package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/peterbourgon/diskv/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Cache struct {
	disabled bool
	d        *diskv.Diskv
}

type Options struct {
	TTL          time.Duration
	ForceRefresh bool
}

type Provider struct {
	Key   any
	Fetch func() (any, error)
}

func (c *Cache) Apply(ctx context.Context, into any, p Provider) error {
	// This probably needs to be checked a bit better, but
	// this is ok for now.
	//
	// We ensure that we can write the data to the pointer/
	// interface that was passed in.
	v := reflect.ValueOf(into).Elem()
	if !v.CanSet() {
		return errors.New("cannot set value here")
	}

	var (
		ttl          time.Duration
		forceRefresh bool
	)

	opts, ok := OptionsFromContext(ctx)
	if ok {
		ttl = opts.TTL
		forceRefresh = opts.ForceRefresh
	}

	if !forceRefresh {
		ok, err := c.Read(p.Key, into, ttl)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	data, err := p.Fetch()
	if err != nil {
		return err
	}

	err = c.Write(p.Key, data)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(data))
	return nil
}

func (c *Cache) Read(key any, into any, ttl time.Duration) (bool, error) {
	if c.disabled {
		return false, nil
	}

	k, err := Key(key)
	if err != nil {
		return false, err
	}

	bytes, err := c.d.Read(k)
	if err != nil {
		log.Warn().Err(err).Msg("")
		return false, nil // FILE MISSING, Do a check here:)
	}

	storedAt, err := Unmarshal(bytes, into)
	if err != nil {
		return false, err
	}

	if storedAt == nil || time.Since(*storedAt) > ttl {
		log.Printf("data expired for %v", key)
		return false, nil
	}

	return true, nil
}

func (c *Cache) Write(key interface{}, data interface{}) error {
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

func NewCache(opts CacheOpts) *Cache {
	if !opts.Enabled {
		return &Cache{disabled: true}
	}

	cacheOpts := diskv.Options{
		BasePath:     opts.BasePath,
		CacheSizeMax: opts.CacheMaxSizeKB * 1024,
	}

	log.Info().Msgf("initializing cache basePath=%s size=%d", cacheOpts.BasePath, cacheOpts.CacheSizeMax)
	return &Cache{d: diskv.New(cacheOpts)}
}

type Value struct {
	At   time.Time
	Data json.RawMessage
}

func Key(data interface{}) (string, error) {
	hash, err := hashstructure.Hash(data, hashstructure.FormatV2, nil)
	if err != nil {
		return "", err
	}

	return strconv.FormatUint(hash, 10), nil
}

func Marshal(data interface{}) ([]byte, error) {
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

func Unmarshal(bytes []byte, into interface{}) (*time.Time, error) {
	var value Value
	if err := json.Unmarshal(bytes, &value); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(value.Data, into); err != nil {
		return nil, err
	}

	return &value.At, nil
}
