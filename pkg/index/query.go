package index

import (
	"os"

	"github.com/bytedance/sonic"
)

func (i *Index) QueryBytes(path string) ([]byte, bool) {
	// Lookup cache
	respBytes, err := i.cache.Get(path)
	if err == nil {
		i.logger.Debugf("cache hit for %s")
		return respBytes, true
	}
	i.logger.Debugf("cache miss for %s: %v", path, err)

	// Query filesystem
	resp, ok := i.queryFromFS(path)
	if !ok {
		return nil, false
	}

	respBytes, err = sonic.Marshal(resp)
	if err != nil {
		i.logger.Errorf("error marshaling response json")
		return nil, false
	}

	// Cache response
	err = i.cache.Set(path, respBytes)
	if err != nil {
		i.logger.Errorf("error saving response to cache")
	}

	return respBytes, false
}

func (i *Index) queryFromFS(path string) (Response, bool) {
	var resp Response
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
			Type:  typeFile,
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

		resp.Type = typeDir
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
				en.Type = typeDir
			} else {
				en.Size = info.Size()
				en.Type = typeFile
			}
			resp.Contents = append(resp.Contents, en)
		}
		return resp, true
	}
}
