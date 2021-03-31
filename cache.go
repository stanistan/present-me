package presentme

import (
	"time"

	c "github.com/stanistan/present-me/internal/cache"
)

var (
	cache    = c.NewCache()
	cacheTTL = 10 * time.Minute
)
