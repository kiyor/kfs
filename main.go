/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Mon 05 Mar 2018 10:57:46 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	// 	"fmt"
	"compress/gzip"
	"flag"
	"github.com/NYTimes/gziphandler"
	// 	quic "github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/h2quic"
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

const (
	KFS_CRT = "/.kfs.crt"
	KFS_KEY = "/.kfs.key"
)

var (
	stop                        bool
	wg                          = new(sync.WaitGroup)
	listen                      = flag.String("l", ":8080", "listen interface")
	listenTLS                   = flag.String("lssl", ":8081", "listen ssl interface")
	rootDir                     = flag.String("root", ".", "root dir")
	gzipTypes                   = flag.String("gzip-types", "text/html text/plain text/css text/javascript text/xml application/json application/javascript application/x-javascript application/xml application/atom+xml application/rss+xml application/vnd.ms-fontobject application/x-font-ttf font/opentype font/x-woff", "gzip type")
	trashPath, crtPath, keyPath string
	enableTLS                   bool
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	trashPath = filepath.Join(*rootDir, "/.Trash")
	if _, err := os.Stat(filepath.Join(*rootDir, "/.Trash")); err != nil {
		os.Mkdir(trashPath, 0744)
	}
	if _, err := os.Stat(filepath.Join(*rootDir, KFS_CRT)); err == nil {
		if _, err = os.Stat(filepath.Join(*rootDir, KFS_KEY)); err == nil {
			enableTLS = true
		}
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

		if req.URL.Path != KFS_CRT && req.URL.Path != KFS_KEY {
			w.Header().Add("Connection", "Keep-Alive")
			w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
			if req.Method == "GET" {
				f := &fileHandler{Dir(*rootDir)}
				f.ServeHTTP(w, req)
			}
		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	mux.Handle("/", handler)
	log.Println("Listen:", *listen, "Root:", *rootDir)
	if enableTLS {
		log.Println("Listen TLS:", *listenTLS, "Root:", *rootDir)
	}
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
		if enableTLS {
			err = h2quic.ListenAndServeQUIC(*listen, crtPath, keyPath, LogHandler(gzipHandler(mux)))
			if err != nil {
				log.Println(err.Error())
				os.Exit(1)
			}
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
