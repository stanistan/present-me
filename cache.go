package presentme

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/mitchellh/hashstructure/v2"
	"github.com/peterbourgon/diskv"
)

var cache = func() *dataCache {
	const (
		cacheDir = "/tmp/crap-data-cache"
	)
	log.Printf("initializing data cache at %s", cacheDir)
	return &dataCache{
		d: diskv.New(diskv.Options{
			BasePath:     cacheDir,
			CacheSizeMax: 10 * 1024,
		}),
	}
}()

type dataCache struct {
	d *diskv.Diskv
}

type cachedModel struct {
	At    time.Time
	Model ReviewModel
}

func (c *dataCache) Read(p *ReviewParams) (*ReviewModel, error) {
	hash, err := hashstructure.Hash(*p, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}

	key := strconv.FormatUint(hash, 10)
	bytes, err := c.d.Read(key)
	if err != nil {
		return nil, nil // FIXME, validate that the file doesn't exist
	}

	var model cachedModel
	if err := json.Unmarshal(bytes, &model); err != nil {
		return nil, err
	}

	if time.Now().Sub(model.At) > 10*time.Minute {
		log.Printf("data expired, will refetch %v", *p)
		return nil, nil
	}

	return &model.Model, nil
}

func (c *dataCache) Write(p *ReviewParams, model *ReviewModel) (*ReviewModel, error) {
	hash, err := hashstructure.Hash(*p, hashstructure.FormatV2, nil)
	if err != nil {
		return nil, err
	}

	key := strconv.FormatUint(hash, 10)
	bytes, err := json.Marshal(cachedModel{At: time.Now(), Model: *model})
	if err != nil {
		return nil, err
	}

	return model, c.d.Write(key, bytes)
}
