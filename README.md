[![Build Status](https://travis-ci.org/yowcow/goromdb.svg?branch=master)](https://travis-ci.org/yowcow/goromdb)

GOROMDB
=======

A single process KVS data store server that talks memcached protocol and serves file-based KVS stores like:

+ JSON file
+ BerkeleyDB file
+ MemcacheDB data in BerkeleyDB file

More protocols and data store may be supported.

HOW TO USE
----------

GOROMDB can be used as an executable binary, or a collection of simple libraries.

### Executable Binary

Just do:

```
go install github.com/yowcow/goromdb
```

To boot:

```
goromdb -addr <address to be bound to> -store <data store> -file <path to data file>
```

An example:

```
goromdb -addr :11211 -store bdb -file path/to/bdb-data.db
```

GOROMDB currently does not daemonize itself.

### Libraries

Just do:

```
go get github.com/yowcow/goromdb
```

and import whatever package into your source code.

BENCHMARK AND PERFORMANCE
-------------------------

GOROMDB should serve fast but maybe not quite as fast as pure memcached.
Detailed benchmark is comming up.

DIRECTORY STRUCTURE
-------------------

When `/tmp/path/to/file.db` is specified to boot option `-file`, GOROMDB creates subdirectories `db00` and `db01` under `/tmp/path/to`,
then start watching for new data file at `/tmp/path/to/file.db` and its checksum file at `/tmp/path/to/file.db.md5`.

```
/tmp/path/to
├── db00
└── db01
```

When data file and its checksum file is placed in directory `/tmp/path/to`, GOROMDB will verify data file against its checksum file.

```
/tmp/path/to
├── db00
├── db01
├── file.db
└── file.db.md5
```

Once checksum succeeds, GOROMDB will move data file into subdirectory either `db00` or `db01`, and load the data into running server.

```
/tmp/path/to
├── db00
│   └── file.db
├── db01
└── file.db.md5
```
