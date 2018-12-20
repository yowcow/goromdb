package main

import (
	"bytes"
	"context"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/yowcow/goromdb/handler/simplehandler"
	"github.com/yowcow/goromdb/loader"
	"github.com/yowcow/goromdb/storage/boltstorage"
	"github.com/yowcow/goromdb/watcher"
)

var (
	concurrency int
	duration    int
	help        bool
	logger      *log.Logger

	bucket      = "goromdb"
	watcherFile = "_watcher/data.db"
	storagePath = "_storage"

	sourceDataFile = "../../data/store/sample-boltdb.db"
	sourceMD5File  = "../../data/store/sample-boltdb.db.md5"
)

func init() {
	flag.IntVar(&concurrency, "c", 1, "concurrency")
	flag.IntVar(&duration, "d", 1, "duration in seconds")
	flag.BoolVar(&help, "h", false, "show help")
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}
}

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func init() {
	if _, err := os.Stat(sourceDataFile); os.IsNotExist(err) {
		logger.Println(sourceDataFile, "is not found")
		os.Exit(1)
	}
	if _, err := os.Stat(sourceMD5File); os.IsNotExist(err) {
		logger.Println(sourceMD5File, "is not found")
		os.Exit(2)
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(duration)*time.Second+3*time.Second)
	defer cancel()

	// create watcher
	wcr := watcher.NewMD5Watcher(watcherFile, 1, logger)
	filein := wcr.Start(ctx)

	// create storage
	stg := boltstorage.New(bucket)

	// create loader
	ldr, err := loader.New(storagePath, "data.db")
	if err != nil {
		panic(err)
	}

	// create handler
	hdr := simplehandler.New(stg, logger)

	var wg sync.WaitGroup

	// start goromdb handler
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		done := hdr.Start(filein, ldr)
		<-done
	}(&wg)

	// start infinite file loading
	wg.Add(1)
	go func(w *sync.WaitGroup) {
		defer w.Done()
		tc := time.NewTicker(500 * time.Millisecond)
		for {
			// copy data.db
			if _, err := os.Stat(watcherFile); os.IsNotExist(err) {
				if r, err := os.Open(sourceDataFile); err == nil {
					if w, err := os.OpenFile(watcherFile+".tmp", os.O_WRONLY|os.O_CREATE, 0644); err == nil {
						io.Copy(w, r)
						w.Close()
					}
					r.Close()
				}
				os.Rename(watcherFile+".tmp", watcherFile)
			}
			// copy data.db.md5
			if _, err := os.Stat(watcherFile + ".md5"); os.IsNotExist(err) {
				if r, err := os.Open(sourceMD5File); err == nil {
					if w, err := os.OpenFile(watcherFile+".md5.tmp", os.O_WRONLY|os.O_CREATE, 0644); err == nil {
						io.Copy(w, r)
						w.Close()
					}
					r.Close()
				}
				os.Rename(watcherFile+".md5.tmp", watcherFile+".md5")
			}

			select {
			case <-tc.C:
			case <-ctx.Done():
				tc.Stop()
				return
			}
		}
	}(&wg)

	time.Sleep(3 * time.Second) // wait 3 secs

	// start infinite `Get` calls
	f := func(id int, w *sync.WaitGroup, c context.Context, l *log.Logger) {
		defer w.Done()
		for {
			_, err := hdr.Get([]byte("hoge"))
			if err != nil {
				l.Println("worker", id, "got error:", err)
			}

			select {
			case <-c.Done():
				return
			default:
			}
		}
	}

	logbuf := new(bytes.Buffer)
	l := log.New(logbuf, "", log.Ldate|log.Ltime)
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go f(i+1, &wg, ctx, l)
	}

	// wait for everybody
	wg.Wait()

	logger.Println("===== errors during `Get()` calls =====")
	io.WriteString(os.Stdout, logbuf.String())
}
