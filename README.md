# urlshort: URL mapping service

Takes the URL GET param h (for hash)
and stores the entire URL into a database.

This database can then be used to forward.

## Example

    localhost:8084/enc/https://www.wikipedia.org?h=wik1
    localhost:8084/enc/https://en.wikipedia.org/wiki/Battle_of_San_Patricio?h=wik2

    http://localhost:8084/dump

    localhost:8084/r/wik1
    localhost:8084/r/wik2