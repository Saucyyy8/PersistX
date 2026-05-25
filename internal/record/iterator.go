package record

import (
	"io"
	"os"
)

// Package record provides an `Iterator` that can read serialized `Record`
// values sequentially from an `*os.File`.
//
// File-level summary:
// - `NewIterator` prepares an `*os.File` for sequential reads (seeks to start).
// - `Iterator.Next` reads the next record by first reading the fixed-size
//   header and then reading the `Key` and `Value` payloads.
//
// Notes on native functions and syntax used below:
// - `(*os.File).Seek(offset int64, whence int)` repositions the file offset.
//   It returns `(int64, error)`. Passing `0, 0` is equivalent to `io.SeekStart`.
// - `io.ReadFull(r, buf)` reads exactly `len(buf)` bytes or returns an error
//   (including `io.EOF` if there aren't enough bytes).
// - `make([]byte, n)` requires `n` to be an `int`; if you have a `uint32`
//   size (like `ks`/`vs`) you need to convert: `make([]byte, int(ks))`.

type Iterator struct {
	file *os.File // underlying file handle
}

// NewIterator returns a new iterator positioned at the start of file `f`.
// Native: `f.Seek(0, 0)` sets the file offset to the start. Seek returns
// `(offset int64, err error)` but here the return values are ignored.
func NewIterator(f *os.File) *Iterator {
	f.Seek(0, 0)
	return &Iterator{file: f}
}

// Next reads the next `Record` from the file.
// Steps:
// 1. Read the fixed-size header using `io.ReadFull` into `headerBuf`.
// 2. `UnmarshalHeader` interprets the header and returns sizes.
// 3. Allocate key/value buffers (remember to convert sizes to `int` when calling `make`).
// 4. Use `io.ReadFull` again to read the exact number of bytes for key/value.
func (it *Iterator) Next() (*Record, error) {
	headerBuf := make([]byte, HeaderSize) // HeaderSize is an untyped constant (int), OK for make()
	_, err := io.ReadFull(it.file, headerBuf)
	if err != nil {
		return nil, err // Could be io.EOF or an actual read error
	}

	ts, t, ks, vs, err := UnmarshalHeader(headerBuf)
	if err != nil {
		return nil, err
	}


	keyBuf := make([]byte, ks) 
	valBuf := make([]byte, vs) 

	if _, err := io.ReadFull(it.file, keyBuf); err != nil {
		return nil, err
	}
	if _, err := io.ReadFull(it.file, valBuf); err != nil {
		return nil, err
	}

	return &Record{
		Timestamp: ts,
		Type:      t,
		KeySize:   ks,
		ValueSize: vs,
		Key:       keyBuf,
		Value:     valBuf,
	}, nil
}
