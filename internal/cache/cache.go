package cache

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/peterbourgon/diskv"
)

type Cache struct {
	d *diskv.Diskv
}

type Provider struct {
	Key          interface{}
	TTL          time.Duration
	ForceRefresh bool
	Fetch        func() (interface{}, error)
}

func (c *Cache) Apply(into interface{}, opts Provider) error {
	// This probably needs to be checked a bit better, but
	// this is ok for now.
	//
	// We ensure that we can write the data to the pointer/
	// interface that was passed in.
	v := reflect.ValueOf(into).Elem()
	if !v.CanSet() {
		return fmt.Errorf("cannot set value here")
	}

	if !opts.ForceRefresh {
		log.Printf("reading from cache")
		ok, err := c.Read(opts.Key, into, opts.TTL)
		if err != nil {
			return err
		}
		if ok {
			return nil
		}
	}

	log.Printf("fetching data")
	data, err := opts.Fetch()
	if err != nil {
		return err
	}

	log.Printf("writing data to cache")
	err = c.Write(opts.Key, data)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(data))
	return nil
}

func (c *Cache) Read(key interface{}, into interface{}, ttl time.Duration) (bool, error) {
	k, err := Key(key)
	if err != nil {
		return false, err
	}

	bytes, err := c.d.Read(k)
	if err != nil {
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

func NewCache() *Cache {
	const cacheDir = "/tmp/present-me-data"
	log.Printf("initializing data cache at %s", cacheDir)
	return &Cache{
		d: diskv.New(diskv.Options{
			BasePath:     cacheDir,
			CacheSizeMax: 10 * 1024,
		}),
	}
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
