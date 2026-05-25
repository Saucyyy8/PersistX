package record

import (
	"encoding/binary"
	"errors"
)

// Package record defines the on-disk binary layout for a single record
// and utilities to marshal/unmarshal that layout.
//
// File-level summary:
// - `Record` is the in-memory representation with fixed-size header fields
//   followed by variable-length `Key` and `Value` byte slices.
// - `Marshal` serializes a `Record` into a single contiguous `[]byte`.
// - `UnmarshalHeader` decodes the fixed-size header so callers know
//   how many bytes to read for `Key` and `Value`.
//
// Notes on native functions and packages used below:
// - `encoding/binary` provides `LittleEndian` helpers:
//   - `PutUint64/PutUint32` write integer values into a byte slice.
//   - `Uint64/Uint32` read integer values from a byte slice.
// - `errors.New` constructs an `error` value from a string.
// - `make([]byte, n)` allocates a byte slice of length `n` (n is `int`).

const (
	// HeaderSize: 8 (Timestamp) + 1 (Type) + 4 (KeySize) + 4 (ValueSize) = 17 bytes
	HeaderSize = 17
	TypePut    = uint8(0)
	TypeDelete = uint8(1)
)

// Record represents the internal "language" of the database.
type Record struct {
	Timestamp int64  // When the record was created
	Type      uint8  // Put or Delete (Tombstone)
	KeySize   uint32 // Length of the key
	ValueSize uint32 // Length of the value
	Key       []byte // Actual key data
	Value     []byte // Actual value data
}

// Marshal converts a Record struct into a raw byte slice.
//
// Steps and native functions used:
// - allocate `buf` using `make([]byte, totalSize)` where `totalSize` is an `int`.
// - `binary.LittleEndian.PutUint64/PutUint32` write integer values into `buf`.
// - `copy(dst, src)` is a built-in that copies bytes from `src` to `dst`.
func (r *Record) Marshal() []byte {
	// 1. Calculate total size to allocate memory once
	totalSize := HeaderSize + len(r.Key) + len(r.Value)
	buf := make([]byte, totalSize) // native `make` for slices

	// 2. Encode fixed-size fields using Little Endian
	// `PutUint64` expects a uint64 value; cast `r.Timestamp` explicitly.
	binary.LittleEndian.PutUint64(buf[0:8], uint64(r.Timestamp))
	buf[8] = r.Type
	binary.LittleEndian.PutUint32(buf[9:13], r.KeySize)
	binary.LittleEndian.PutUint32(buf[13:17], r.ValueSize)

	// 3. Copy variable-length data using the built-in `copy`
	copy(buf[HeaderSize:HeaderSize+int(r.KeySize)], r.Key)
	copy(buf[HeaderSize+int(r.KeySize):], r.Value)

	return buf
}

// UnmarshalHeader reads the first 17 bytes to understand how much data follows.
//
// Native: `binary.LittleEndian.Uint64/Uint32` read unsigned integers from
// a byte slice using little-endian byte order.
// Returns an error constructed with `errors.New` when `data` is too short.
func UnmarshalHeader(data []byte) (int64, uint8, uint32, uint32, error) {
	if len(data) < HeaderSize {
		return 0, 0, 0, 0, errors.New("insufficient data for header")
	}
	ts := int64(binary.LittleEndian.Uint64(data[0:8]))
	t := data[8]
	ks := binary.LittleEndian.Uint32(data[9:13])
	vs := binary.LittleEndian.Uint32(data[13:17])
	return ts, t, ks, vs, nil
}
