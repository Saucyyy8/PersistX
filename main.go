package main

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Saucyyy8/persistx/internal/filestore"
	"github.com/Saucyyy8/persistx/internal/record"
)

// main is a small demonstration of writing and reading a `record.Record`
// using the `filestore` and `record` packages.
//
// File-level summary:
//   - Open or create an append-only log file with `filestore.NewFileStore`.
//   - Construct a `record.Record`, marshal it to bytes and append it to the file.
//   - Seek back to the file start, read the header with `io.ReadFull`,
//     unmarshal with `record.UnmarshalHeader`, then read the key/value payloads.
//
// Notes on native functions and conversions used below:
// - `time.Now().Unix()` returns an `int64` (seconds since epoch).
// - `[]byte(string)` converts a Go `string` to a `[]byte` copy.
// - `string([]byte)` converts a `[]byte` back to a `string`.
// - `defer store.Close()` calls `(*os.File).Close()` when `main` exits.
// - `log.Fatalf` prints a formatted message and then calls `os.Exit(1)`.
// - `io.ReadFull(r, buf)` reads exactly `len(buf)` bytes or returns an error.
// - `(*os.File).Seek(offset int64, whence int)` repositions the file offset.

func RecoveryEngine(store *filestore.FileStore) map[string]string {
	state := make(map[string]string)
	iter := record.NewIterator(store.Handle)

	for {
		rec, err := iter.Next()
		if err == io.EOF {
			break // Reached the end of the log
		}
		if err != nil {
			log.Printf("Warning: Corrupt record found during recovery: %v", err)
			break
		}

		// Apply the log entry to our state
		key := string(rec.Key)
		if rec.Type == record.TypeDelete {
			delete(state, key)
		} else {
			state[key] = string(rec.Value)
		}
	}
	return state
}
func main() {
	// Initialize our storage
	store, err := filestore.NewFileStore("persistx.log") // native: os.OpenFile inside
	if err != nil {
		log.Fatalf("Failed to open store: %v", err) // native: log.Fatalf -> prints then exits
	}
	defer store.Close() // native: Close releases file descriptor when main returns

	// 1. Define a Record
	key := "user:100"
	val := "Pranav"

	r := &record.Record{
		Timestamp: time.Now().Unix(), // native: returns int64
		Type:      record.TypePut,
		KeySize:   uint32(len(key)), // convert int -> uint32 for header
		ValueSize: uint32(len(val)),
		Key:       []byte(key), // native: string -> []byte conversion (makes a copy)
		Value:     []byte(val),
	}

	store.Write(r.Marshal())
	store.Sync()

	fmt.Println("Intital write synced to disk.")

	store.Close()
	fmt.Println("-- System Crashed Simulation --")
	newStore, _ := filestore.NewFileStore("persistx.log")
	defer newStore.Close()

	dbState := RecoveryEngine(newStore)

	fmt.Printf("Recovered State: %+v\n", dbState)
	
	if dbState["user:101"] == "Suhas" {
		fmt.Println("SUCCESS: Data survived the crash!")
	} else {
		fmt.Println("FAILURE: Data lost.")
	}
	fmt.Println(dbState)
}
