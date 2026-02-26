package index

import (
	"fmt"
	"time"

	"github.com/allegro/bigcache"
)

type Index struct {
	cache  *bigcache.BigCache
	logger Logger

	// Config
	root    string
	ttl     time.Duration
	maxSize int64
}

func New(opts ...func(*Index)) (*Index, error) {
	index := &Index{
		root:    ".",
		ttl:     time.Minute,
		maxSize: 10,
		logger:  &DiscardLogger{},
	}
	for _, o := range opts {
		o(index)
	}
	var err error
	index.cache, err = bigcache.NewBigCache(bigcache.Config{
		Shards:     1024,
		LifeWindow: index.ttl,
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

func WithMaxSize(size int64) func(*Index) {
	return func(i *Index) {
		i.maxSize = size
	}
}

func WithLogger(logger Logger) func(*Index) {
	return func(i *Index) {
		i.logger = logger
	}
}

func (i *Index) Close() error {
	return i.cache.Close()
}
