# PersistX

Work in progress: I am still actively building and refining this project.

## Overview
PersistX is a full database storage engine built from scratch in Go, modeled after LevelDB/RocksDB (LSM-Tree architecture) and extended with AI/ML features. It focuses on append-only durability, immutable on-disk structures, and fast reads through multi-level indexing, while keeping the write path minimal and crash-safe.

## What It Does
- Provides a storage engine with a write-ahead log, in-memory indexing, immutable SSTables, and a multi-level read path.
- Uses background flushing and compaction to keep on-disk data organized and efficient.
- Adds AI/ML layers for semantic search, learned filters, adaptive caching, and asynchronous intelligence pipelines.
- Exposes a TCP interface for basic commands and client access.

## Core Storage Engine
- Binary I/O and record serialization using fixed metadata fields and byte-level encoding.
- Append-only WAL with fsync durability and crash recovery via log replay.
- Thread-safe in-memory index (Skip-List style) with size tracking and freeze thresholds.
- Background flushing to immutable SSTables with index/footer for fast seeks.
- Read hierarchy across memtables and SSTables with Bloom filters to avoid unnecessary disk I/O.
- Compaction via k-way merge, with tombstone handling to clean up deleted keys.
- TCP server interface with simple text commands for SET, GET, and DEL.

## AI/ML Extensions
- Semantic layer that generates embeddings on writes and supports vector-based search via HNSW.
- Learned Bloom filters trained on access patterns to reduce disk reads, with a standard Bloom fallback to avoid false negatives.
- Adaptive cache policies using LeCaR to balance LRU and LFU behavior with reinforcement signals.
- Asynchronous worker pipelines for embeddings, inference, and index updates with backpressure controls.
- Benchmarking and evaluation to measure latency, disk savings, cache hit ratios, and semantic search quality.

## Tech Stack
- Language: Go
- Storage model: LSM-Tree
- Inspiration: LevelDB, RocksDB
- AI models: all-MiniLM-L6-v2 embeddings, learned Bloom classifier
- Concurrency: goroutines, sync.RWMutex, worker pools
- Networking: Go net package (TCP)

## Design Principles
- Append-only writes (WAL and SSTables)
- Immutable on-disk structures
- Async AI processing that never blocks the write path
- Crash safety via WAL replay
- Multi-level read path with probabilistic and learned filters
