/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Fri 05 Jan 2018 01:01:10 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"fmt"
	"compress/gzip"
	"flag"
	"github.com/NYTimes/gziphandler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	stop      bool
	wg        = new(sync.WaitGroup)
	listen    = flag.String("l", ":8080", "listen interface")
	rootDir   = flag.String("root", ".", "root dir")
	gzipTypes = flag.String("gzip-types", "text/plain text/css text/javascript text/xml application/json application/javascript application/x-javascript application/xml application/atom+xml application/rss+xml application/vnd.ms-fontobject application/x-font-ttf font/opentype font/x-woff", "gzip type")
	trash     string
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	trash = filepath.Join(*rootDir, "/.Trash")
	if _, err := os.Stat(filepath.Join(*rootDir, "/.Trash")); err != nil {
		os.Mkdir(trash, 0744)
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	mux := http.NewServeMux()
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if stop {
			return
		}
		wg.Add(1)
		defer wg.Done()

		w.Header().Add("Connection", "Keep-Alive")
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		if req.Method == "GET" {
			f := &fileHandler{Dir(*rootDir)}
			f.ServeHTTP(w, req)
			// 		} else if req.Method == "POST" || req.Method == "PUT" {
			// 			uploadHandler(w, req)
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	mux.Handle("/", handler)
	log.Println("Listen:", *listen, "Root:", *rootDir)
	go func() {
		typ := strings.Split(*gzipTypes, " ")
		for _, v := range typ {
			typ = append(typ, v+"; charset=utf-8")
		}
		gzipHandler, err := gziphandler.GzipHandlerWithOpts(gziphandler.MinSize(512), gziphandler.CompressionLevel(gzip.DefaultCompression), gziphandler.ContentTypes(typ))
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
		err = http.ListenAndServe(*listen, LogHandler(gzipHandler(mux)))
		if err != nil {
			log.Println(err.Error())
			os.Exit(1)
		}
	}()

forever:
	for {
		select {
		case s := <-sig:
			log.Printf("Signal (%d) received, stopping\n", s)
			stop = true
			log.Println("stopped")
			break forever
		}
	}
}
