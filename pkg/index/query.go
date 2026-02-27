package index

import (
	"encoding/binary"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bytedance/sonic"
)

func (i *Index) Query(path string) (Response, bool) {
	var resp Response
	respBytes, ok := i.QueryBytes(path)
	if !ok {
		return resp, false
	}
	err := sonic.Unmarshal(respBytes, &resp)
	if err != nil {
		i.logger.Errorf("response unmarshal failed: %v", err)
		return resp, false
	}
	return resp, true
}

func (i *Index) QueryBytes(path string) ([]byte, bool) {
	// Strip trailing slash to avoid duplicate cache
	path = strings.TrimSuffix(path, "/")
	i.logger.Debugf("query \"%s\"", path)

	// Lookup cache
	respBytes, ok := i.queryCache(path)
	if ok {
		return respBytes, true
	}

	// Query filesystem
	resp, ok := i.queryFilesystem(path)
	if !ok {
		i.logger.Debugf("not found on filesystem: %s", path)
		return nil, false
	}

	respBytes, err := sonic.Marshal(resp)
	if err != nil {
		i.logger.Errorf("error marshaling response json")
		return nil, false
	}

	// Cache response
	err = i.putCache(path, respBytes)
	if err != nil {
		i.logger.Errorf("error saving response to cache")
	}

	return respBytes, true
}

type cacheHeader struct {
	ExpiresAt int64 // Unix timestamp
}

// Use special header to handle expiry
func (i *Index) queryCache(path string) ([]byte, bool) {
	respBytes, err := i.cache.Get(path)
	if err == nil {
		header, body, err := extractHeader(respBytes)
		if err != nil {
			i.logger.Errorf("error extracting header: %v", err)
			return nil, false
		}
		if time.Now().Unix() >= header.ExpiresAt {
			i.logger.Debugf("cache expired for \"%s\"", path)
			return nil, false
		}
		i.logger.Debugf("cache hit for \"%s\"", path)
		return body, true
	}
	i.logger.Debugf("cache miss for \"%s\"", path)
	return nil, false
}

func (i *Index) putCache(path string, respBytes []byte) error {
	data := prependHeader(respBytes, cacheHeader{
		ExpiresAt: time.Now().Add(i.ttl).Unix(),
	})
	return i.cache.Set(path, data)
}

func extractHeader(data []byte) (cacheHeader, []byte, error) {
	if len(data) < 8 {
		return cacheHeader{}, nil, errors.New("not enough bytes")
	}
	header := cacheHeader{
		ExpiresAt: int64(binary.BigEndian.Uint64(data[:8])),
	}
	return header, data[8:], nil
}

func prependHeader(body []byte, header cacheHeader) []byte {
	buf := make([]byte, 8+len(body))
	binary.BigEndian.PutUint64(buf[:8], uint64(header.ExpiresAt))
	copy(buf[8:], body)
	return buf
}

func (i *Index) queryFilesystem(path string) (Response, bool) {
	var resp Response
	path = filepath.Join(i.root, path)

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			i.logger.Debugf("path %s not found", path)
			return resp, false
		}
		i.logger.Errorf("error opening path %s: %v", path, err)
		return resp, false
	}

	if !info.IsDir() {
		// Handle file
		return Response{
			Type:  TypeFile,
			MTime: info.ModTime().Unix(),
			Size:  info.Size(),
		}, true
	} else {
		// Handle directory
		entries, err := os.ReadDir(path)
		if err != nil {
			i.logger.Errorf("error reading directory %s: %v", path, err)
			return resp, false
		}

		resp.Type = TypeDir
		resp.Contents = make([]Entry, 0, len(entries))

		for _, e := range entries {
			info, err := e.Info()
			if err != nil {
				i.logger.Warnf("error getting info of entry %s/%s: %v", path, e.Name(), err)
				continue
			}
			en := Entry{
				Name:  info.Name(),
				MTime: info.ModTime().Unix(),
			}
			if info.IsDir() {
				en.Type = TypeDir
			} else {
				en.Size = info.Size()
				en.Type = TypeFile
			}
			resp.Contents = append(resp.Contents, en)
		}
		return resp, true
	}
}
