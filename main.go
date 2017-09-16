/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Thu 07 Sep 2017 01:07:06 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"fmt"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"syscall"
)

var (
	stop    bool
	wg      = new(sync.WaitGroup)
	listen  = flag.String("l", ":8080", "listen interface")
	rootDir = flag.String("root", ".", "root dir")
	trash   string
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

	mux.Handle("/", LogHandler(handler))
	log.Println("Listen:", *listen, "Root:", *rootDir)
	go func() {
		err := http.ListenAndServe(*listen, mux)
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
