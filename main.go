/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Tue Aug 29 12:44:47 2017

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
	"runtime"
	"sync"
	"syscall"
)

var (
	stop bool
	wg   = new(sync.WaitGroup)
	// 	ch     = make(chan bool)
	listen = flag.String("l", ":8080", "listen interface")
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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
		// 		ch <- true

		w.Header().Add("Connection", "Keep-Alive")
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		if req.Method == "GET" {
			f := &fileHandler{Dir(".")}
			f.ServeHTTP(w, req)
			// 		} else if req.Method == "POST" || req.Method == "PUT" {
			// 			uploadHandler(w, req)
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	mux.Handle("/", LogHandler(handler))
	log.Println("Listen:", *listen)
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
