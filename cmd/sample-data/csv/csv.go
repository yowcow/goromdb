package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"os"
)

type Data map[string]string

func main() {
	var jsonFile string
	var csvFile string

	flag.StringVar(&jsonFile, "input-from", "data/sample-data.json", "read JSON from")
	flag.StringVar(&csvFile, "output-to", "data/sample-radix.csv", "write CSV to")
	flag.Parse()

	writeCSV(jsonFile, csvFile)
}

func writeCSV(jsonFile, csvFile string) {
	fi, err := os.Open(jsonFile)
	if err != nil {
		panic(err)
	}
	defer fi.Close()

	var data Data
	decoder := json.NewDecoder(fi)
	if err = decoder.Decode(&data); err != nil {
		panic(err)
	}

	fo, err := os.Create(csvFile)
	if err != nil {
		panic(err)
	}
	defer fo.Close()

	w := csv.NewWriter(fo)
	err = w.Write([]string{"key", "value"}) // header line
	if err != nil {
		panic(err)
	}
	w.Flush()

	for k, v := range data {
		err = w.Write([]string{k, v})
		if err != nil {
			panic(err)
		}
		w.Flush()
	}
}
