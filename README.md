# urlshort: URL mapping service

Takes the URL GET param h (for hash)  
and stores the URL under this key into a database.

The key can then be used to retrieve the URL  
from the database.
Or to be forwarded / redirected to the URL.

## Storage Engines

The filesystem storage is from someone else.
The implementation is counting all existing files on each insert.
Dont use it.

The BoltDB storage is not fully implemented.
BoltDB is supposed to be slightly better at retrieval than at insert.

The LevelDB storage is workable and yields some encouraging benchmark results.

## Example

### Store/Save

    localhost:8084/enc/https://www.wikipedia.org?h=wik1
    localhost:8084/enc/https://en.wikipedia.org/wiki/Battle_of_San_Patricio?h=wik2

### Check all

    http://localhost:8084/dump

### Check one

    http://localhost:8084/dec/wik1

### Call to redirect

    http://localhost:8084/r/wik1
    http://localhost:8084/r/wik2

## Benchmark LevelDB

8 microseconds per saving operation
125.000 inserts per second.

5 microseconds per loading operation
200.000 loads per second.

