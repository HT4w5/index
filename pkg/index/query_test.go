package index

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/HT4w5/autoindex/pkg/log"
)

const (
	maxNameLen  = 32
	maxFileSize = 1024
	charSet     = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func randomName(r *rand.Rand) string {
	l := r.IntN(maxNameLen) + 1
	var buf bytes.Buffer
	for range l {
		buf.WriteByte(charSet[r.IntN(len(charSet))])
	}
	return buf.String()
}

func makeBenchmarkDir(b *testing.B, r *rand.Rand, nFiles int, nDirs int) string {
	dir := b.TempDir()

	// Create files
	for range nFiles {
		f, err := os.Create(filepath.Join(dir, randomName(r)))
		if err != nil {
			b.Fatalf("create error: %v", err)
		}
		err = f.Truncate(r.Int64N(maxFileSize))
		if err != nil {
			b.Fatalf("truncate error: %v", err)
		}
	}

	// Create directories
	for range nDirs {
		d := filepath.Join(dir, randomName(r))
		err := os.Mkdir(d, 0600)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				continue
			} else {
				b.Fatalf("mkdir error: %v", err)
			}
		}
	}

	return dir
}

const (
	nFiles = 1000
	nDirs  = 1000
	seed   = 1024
)

func BenchmarkQueryFilesystem(b *testing.B) {
	var seedBytes [32]byte
	binary.BigEndian.PutUint64(seedBytes[:], seed)
	r := rand.New(rand.NewChaCha8(seedBytes))
	dir := makeBenchmarkDir(b, r, nFiles, nDirs)
	idx := Index{
		root:   dir,
		logger: &log.DiscardLogger{},
	}

	b.ResetTimer()

	for range b.N {
		_, ok := idx.queryFilesystem("")
		if !ok {
			b.Fatal("query failed")
		}
	}
}

func BenchmarkQueryCache(b *testing.B) {
	var seedBytes [32]byte
	binary.BigEndian.PutUint64(seedBytes[:], seed)
	r := rand.New(rand.NewChaCha8(seedBytes))
	dir := makeBenchmarkDir(b, r, nFiles, nDirs)
	idx, err := New(
		WithLogger(&log.DiscardLogger{}),
		WithRoot(dir),
		WithTTL(10*time.Minute),
	)
	if err != nil {
		b.Fatalf("failed to create index: %v", err)
	}

	_, ok := idx.QueryBytes("/")
	if !ok {
		b.Fatal("query failed")
	}

	b.ResetTimer()

	for range b.N {
		_, ok := idx.QueryBytes("/")
		if !ok {
			b.Fatal("query failed")
		}
	}
}
