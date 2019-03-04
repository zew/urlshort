# urlshort: URL mapping service

Takes the URL GET param h (for hash)  
and stores the URL under this key into a database.

The key can then be used to retrieve the URL  
from the database.
Or to forward/redirect to the URL.

## Storage Engines

The LevelDB and BoltDB storages are workable and yield interesting benchmark results.  

BoltDB is supposed to be slightly better at retrieval than at insert.  
BoltDB's insertion speed seems relatively bad.
That's because it's flushed/synced after each insert.

LevelDB is flushed/synced after an unclear number of inserts.  
If we flush/sync on _each_ write, performance breaks down too.

This is expected behavior - compare the C implementation docs:  
https://github.com/google/leveldb/blob/master/doc/index.md

Also the internal implementation notes:  
https://github.com/google/leveldb/blob/master/doc/impl.md

LevelDB is susceptible to corruption.  
Use Recover() when corrupted (not implemented).

A parallel logfile (search code for `transaction log`) was coded, but discarded.  
There are levelDB journal files. They can be processed by https://godoc.org/github.com/syndtr/goleveldb/leveldb/journal.

    The wire format allows for limited recovery in the face of data corruption:  
    on a format error (such as a checksum mismatch), the reader moves to the next block  
    and looks for the next full or first chunk.

We cleanly close databases after HTTP server crash or regular application termination.

## Example

### Store/Save

    localhost:8084/enc?url=https://www.wikipedia.org?h=wik1
    localhost:8084/enc?url=https://en.wikipedia.org/wiki/Battle_of_San_Patricio?h=wik2

### Check all

    http://localhost:8084/dump

### Check one

    http://localhost:8084/dec/wik1

### Call to redirect

    http://localhost:8084/r/wik1
    http://localhost:8084/r/wik2

## Benchmark LevelDB

8-11 microseconds per saving operation -  
 >90.000 inserts per second.

3-5 microseconds per loading operation -  
 >200.000 loads per second.

### C implementation

Source: https://github.com/google/leveldb

    fillseq      :       1.765 micros/op;   62.7 MB/s
    fillsync     :     268.409 micros/op;    0.4 MB/s (10000 ops)
    fillrandom   :       2.460 micros/op;   45.0 MB/s
    overwrite    :       2.380 micros/op;   46.5 MB/s

## Benchmark BoltDB

2.500 microseconds per saving operation -  
 400 inserts per second.

2 microseconds per loading operation -  
 500.000 loads per second.

## Database considerations

Possible database backends - LevelDB vs BoldDB

LevelDB and its derivatives (RocksDB, HyperLevelDB) underlying structure is a log-structured merge-tree (LSM tree).  
 An LSM tree optimizes random writes by using a write ahead log and multi-tiered, sorted files called SSTables.  
 Bolt uses a B+tree internally and only a single file.

https://github.com/etcd-io/bbolt

https://github.com/syndtr/goleveldb

### Other

https://github.com/tidwall/buntdb based on https://github.com/tidwall/btree

https://github.com/dgraph-io/badger
