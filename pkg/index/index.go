package index

import (
	"fmt"
	"time"

	"github.com/HT4w5/autoindex/pkg/log"
	"github.com/allegro/bigcache"
	"github.com/docker/go-units"
)

type Index struct {
	cache  *bigcache.BigCache
	logger log.Logger

	// Config
	root    string
	ttl     time.Duration
	maxSize int
}

func New(opts ...func(*Index)) (*Index, error) {
	index := &Index{
		root:    ".",
		ttl:     time.Minute,
		maxSize: 10,
		logger:  &log.DiscardLogger{},
	}
	for _, o := range opts {
		o(index)
	}
	var err error
	index.cache, err = bigcache.NewBigCache(bigcache.Config{
		Shards:             1024,
		LifeWindow:         index.ttl,
		MaxEntriesInWindow: max(100, index.maxSize*units.MB/(10*units.KB)),
		MaxEntrySize:       10 * units.KB,
		CleanWindow:        time.Minute,
		HardMaxCacheSize:   index.maxSize,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating bigcache: %w", err)
	}
	return index, nil
}

func WithRoot(root string) func(*Index) {
	return func(i *Index) {
		i.root = root
	}
}

func WithTTL(ttl time.Duration) func(*Index) {
	return func(i *Index) {
		i.ttl = ttl
	}
}

func WithMaxSize(size int) func(*Index) {
	return func(i *Index) {
		i.maxSize = size
	}
}

func WithLogger(logger log.Logger) func(*Index) {
	return func(i *Index) {
		i.logger = logger
	}
}

func (i *Index) Close() error {
	return i.cache.Close()
}
