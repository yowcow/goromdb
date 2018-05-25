package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"

	"github.com/boltdb/bolt"
)

// Data represents a key-value data
type Data map[string]string

// NSData represents a namespaced Data
type NSData map[string]Data

func main() {
	var jsonFile string
	var dbFile string
	var dbBucket string

	flag.StringVar(&jsonFile, "input-from", "data/sample-ns-data.json", "read JSON from")
	flag.StringVar(&dbFile, "output-to", "data/sample-ns-boltdb.db", "write database to")
	flag.Parse()

	writeDB(jsonFile, dbFile, dbBucket)
}

func writeDB(jsonFile, dbFile, dbBucket string) {
	var nsdata NSData

	b, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(b, &nsdata)
	if err != nil {
		panic(err)
	}

	db, err := bolt.Open(dbFile, 0644, nil)
	if err != nil {
		panic(err)
	}

	for bucket, data := range nsdata {
		// create buckets
		db.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucket([]byte(bucket))
			if err != nil {
				panic(err)
			}
			return nil
		})

		// insert keys
		for k, v := range data {
			db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucket))
				err := b.Put([]byte(k), []byte(v))
				return err
			})
		}
	}
}
