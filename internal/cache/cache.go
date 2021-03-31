package cache

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/peterbourgon/diskv"
)

type Cache struct {
	d *diskv.Diskv
}

func (c *Cache) Read(key interface{}, into interface{}, ttl time.Duration) error {
	k, err := Key(key)
	if err != nil {
		return err
	}

	bytes, err := c.d.Read(k)
	if err != nil {
		return nil // FILE MISSING, Do a check here:)
	}

	storedAt, err := Unmarshal(bytes, into)
	if err != nil {
		return err
	}

	if storedAt == nil || time.Now().Sub(*storedAt) > ttl {
		log.Printf("data expired for %v", key)
		into = nil
		return nil
	}

	return nil
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
