package integrity

// Manifest represents the structure of the server-side modpack manifest
type Manifest struct {
	Version int        `json:"version"`
	Files   []FileInfo `json:"files"`
}

// FileInfo represents a single file trackable by the integrity system
type FileInfo struct {
	Path     string `json:"path"`
	Size     int64  `json:"size"`
	Checksum string `json:"checksum"`
	Override bool   `json:"override"`
}
