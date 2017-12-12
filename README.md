[![Build Status](https://travis-ci.org/yowcow/goromdb.svg?branch=master)](https://travis-ci.org/yowcow/goromdb)

GOROMDB
=======

Yet another single process KVS server implemented over file-based database.

GOROMDB is a datastore that:

+ accepts read-only access
+ talks multiple protocols like memcached and others
+ handles multiple database backend like JSON, BerkeleyDB, BoltDB, and others

HOW TO USE
----------

GOROMDB can be used as an executable binary, or a collection of simple libraries.

### Executable Binary

Just do:

```
go install github.com/yowcow/goromdb
```

An example of booting GOROMDB with radix index tree and BerkeleyDB database is:

```
goromdb -addr :11211 -handler radix -storage bdb -file path/to/bdb-data.db -basedir path/to/store
```

GOROMDB does not daemonize itself.

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

When `-basedir` of `/tmp/path/to/dir` is specified at boot, GOROMDB creates subdirectories `data00` and `data01` under `/tmp/path/to/dir`:

```
/tmp/path/to/dir
├── db00
└── db01
```

When `-file` of `/tmp/path/to/dir/data.db` is specified at boot, GOROMDB will watch for database file `/tmp/path/to/dir/data.db` and its MD5 sum file `/tmp/path/to/dir/data.db.md5`.

When database and MD5 files are placed in directory `/tmp/path/to/dir`, GOROMDB will verify MD5 sum.

```
/tmp/path/to/dir
├── data00
├── data01
├── data.db
└── data.db.md5
```

Once MD5 sum verification succeeds, GOROMDB will move data file into subdirectory either `data00` or `data01`, and load the database into running server.

```
/tmp/path/to/dir
├── data00
│   └── data.db
└── data01
```

Placing database and MD5 sum files again will load database into next subdirectory `data01` vice versa.
