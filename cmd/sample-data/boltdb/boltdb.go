package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/boltdb/bolt"
)

// Data represents a key-value data
type Data map[string]string

func main() {
	var jsonFile string
	var dbFile string
	var dbBucket string

	flag.StringVar(&jsonFile, "input-from", "data/sample-data.json", "read JSON from")
	flag.StringVar(&dbFile, "output-to", "data/sample-boltdb.db", "write database to")
	flag.StringVar(&dbBucket, "output-bucket", "goromdb", "bucket name to put data")
	flag.Parse()

	writeDB(jsonFile, dbFile, dbBucket)
}

func writeDB(jsonFile, dbFile, dbBucket string) {
	var data Data

	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open(dbFile, 0644, nil)
	if err != nil {
		panic(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucket([]byte(dbBucket))
		if err != nil {
			panic(err)
		}
		return nil
	})

	for k, v := range data {
		db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(dbBucket))
			err := b.Put([]byte(k), []byte(v))
			return err
		})
	}
}
