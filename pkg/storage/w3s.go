package storage

import (
	"bytes"
	"io/fs"
	"time"
)

// IntermediateFile implements fs.File and is used to wrap the data read from GCS.
// This is needed because the Put method of the w3s.Client interface expects a fs.File.
// The Put method is used to upload the data to web3.storage.
type IntermediateFile struct {
	data *bytes.Reader
	name string
}

// NewIntermediateFile creates a new IntermediateFile instance.
func NewIntermediateFile(data []byte, name string) *IntermediateFile {
	return &IntermediateFile{
		data: bytes.NewReader(data),
		name: name,
	}
}

// Implement fs.File methods //

// Stat returns a fs.FileInfo describing the file.
func (f *IntermediateFile) Stat() (fs.FileInfo, error) {
	return &fileInfo{name: f.name, size: f.data.Size()}, nil
}

// Read reads up to len(p) bytes into p.
func (f *IntermediateFile) Read(p []byte) (n int, err error) {
	return f.data.Read(p)
}

// Close noop for an in memory virtual file.
func (f *IntermediateFile) Close() error {
	return nil
}

type fileInfo struct {
	name string
	size int64
}

// Implement fs.FileInfo methods //

func (fi *fileInfo) Name() string       { return fi.name }
func (fi *fileInfo) Size() int64        { return fi.size }
func (fi *fileInfo) Mode() fs.FileMode  { return 0o444 } // Read-only
func (fi *fileInfo) ModTime() time.Time { return time.Now() }
func (fi *fileInfo) IsDir() bool        { return false }
func (fi *fileInfo) Sys() interface{}   { return nil }
