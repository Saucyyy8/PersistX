package filestore

import (
	"os"
)

// Package filestore provides a thin wrapper around an *os.File
// to support simple append-only storage for records.
//
// File-level summary:
// - `FileStore` holds an `*os.File` handle used for all file operations.
// - `NewFileStore` opens (or creates) the underlying file with append
//   semantics and returns a `*FileStore`.
// - `Write` appends a byte slice to the file.
// - `Sync` flushes in-memory file buffers to stable storage.
// - `Close` closes the file descriptor.
//
// Note on native functions and values used below:
// - `os.OpenFile(path, flags, perm)` opens or creates a file and returns
//   `(*os.File, error)`.
// - `os.File.Write([]byte)` writes bytes to the file and returns `(int, error)`.
// - `os.File.Sync()` forces the OS to flush file buffers to disk.
// - `os.File.Close()` closes the file descriptor and releases OS resources.
// - Flag constants like `os.O_APPEND`, `os.O_CREATE`, `os.O_RDWR` are used
//   together with bitwise OR to configure open behavior.
// - File permission `0644` is an octal literal (owner read/write, group/world read).

type FileStore struct {
	Handle *os.File // native type `*os.File` from package `os` representing an open file.
}

// Sync calls the underlying `(*os.File).Sync()` method.
// `Sync` flushes the file's in-memory copy to stable storage (fsync).
func (fs *FileStore) Sync() error {
	return fs.Handle.Sync()
}

// NewFileStore opens or creates a log file in append-only mode and
// returns a `*FileStore`.
//
// Native: `os.OpenFile` - signature `func OpenFile(name string, flag int, perm FileMode) (*File, error)`.
// The provided flags are combined using bitwise OR:
// - `os.O_APPEND`: always write to the end of the file.
// - `os.O_CREATE`: create the file when it does not exist.
// - `os.O_RDWR`: open for read and write.
// `perm` (0644) is an octal integer literal describing file permissions.
func NewFileStore(path string) (*FileStore, error) {
	// open (or create) file with append and read/write semantics
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	return &FileStore{Handle: f}, nil
}

// Write appends raw bytes to the disk.
//
// Native: `(*os.File).Write` writes `data` to the file and returns `(n int, err error)`.
// We discard the written byte count and propagate the `error`.
func (fs *FileStore) Write(data []byte) error {
	_, err := fs.Handle.Write(data)
	return err
}

// Close closes the underlying file descriptor.
// Native: `(*os.File).Close()` releases OS resources and may return an error.
func (fs *FileStore) Close() error {
	return fs.Handle.Close()
}
