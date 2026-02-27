package index

const (
	TypeFile = "file"
	TypeDir  = "dir"
)

type Response struct {
	Type     string  `json:"type"` // "file" or "dir"
	MTime    int64   `json:"mtime,omitempty"`
	Size     int64   `json:"size,omitempty"`
	Contents []Entry `json:"content,omitempty"`
}

type Entry struct {
	Name  string `json:"name"`
	Type  string `json:"type"`  // "file" or "dir"
	MTime int64  `json:"mtime"` // Unix timestamp
	Size  int64  `json:"size,omitempty"`
}
