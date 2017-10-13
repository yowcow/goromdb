package main

import (
	"flag"

	"github.com/ajiyoshi-vg/goberkeleydb/bdb"
)

type Data map[string]string

func main() {
	var file string

	flag.StringVar(&file, "output-to", "data/sample-bdb.db", "write database to")
	flag.Parse()

	writeDB(file)
	readDB(file)
}

func readDB(file string) {
	db, err := bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbReadOnly, 0)

	if err != nil {
		panic(err)
	}

	keys := []string{
		"hoge",
		"fuga",
		"foo",
		"bar",
		"buz",
	}

	for _, k := range keys {
		if _, err := db.Get(bdb.NoTxn, []byte(k), 0); err != nil {
			panic(err)
		}
	}
}

func writeDB(file string) {
	db, err := bdb.OpenBDB(bdb.NoEnv, bdb.NoTxn, file, nil, bdb.BTree, bdb.DbCreate, 0)

	if err != nil {
		panic(err)
	}

	data := Data{
		"hoge": "hoge!",
		"fuga": "fuga!!",
		"foo":  "foo!!!",
		"bar":  "bar!!!!",
		"buz":  "buz!!!!!",
	}

	for k, v := range data {
		if err = db.Put(bdb.NoTxn, []byte(k), []byte(v), 0); err != nil {
			panic(err)
		}
	}

	db.Close(0)
}
