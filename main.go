/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : main.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Mon 21 May 2018 04:19:32 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"compress/gzip"
	"flag"
	"github.com/NYTimes/gziphandler"
	// 	"github.com/aws/aws-sdk-go/aws"
	awscredentials "github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	minio "github.com/minio/minio-go"
	"github.com/minio/minio-go/pkg/credentials"
	"log"
	"net/http"
	"os"
	"os/signal"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"syscall"
)

var (
	stop                        bool
	wg                          = new(sync.WaitGroup)
	listen                      = flag.String("l", ":8080", "listen interface")
	rootDir                     = flag.String("root", ".", "root dir")
	gzipTypes                   = flag.String("gzip-types", "text/html text/plain text/css text/javascript text/xml application/json application/javascript application/x-javascript application/xml application/atom+xml application/rss+xml application/vnd.ms-fontobject application/x-font-ttf font/opentype font/x-woff", "gzip type")
	trashPath, crtPath, keyPath string

	s3proxy    = flag.Bool("s3proxy", false, "use as s3 proxy, require key")
	awsProfile = flag.String("profile", "dos", "aws profile store the key")
	s3endpoint = flag.String("s3endpoint", "sfo2.digitaloceanspaces.com", "s3 endpoint")
	s3region   = flag.String("region", "sfo2", "s3 region")
	s3bucket   = flag.String("s3bucket", "dos", "s3 bucket")

	s3client *minio.Client

	signer *v4.Signer
)

func init() {
	flag.Parse()
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if *s3proxy {
		cuser, err := user.Current()
		if err != nil {
			cuser, _ = user.LookupId("0")
		}
		authfile := filepath.Join(cuser.HomeDir, ".aws/credentials")
		credential := credentials.NewFileAWSCredentials(authfile, *awsProfile)
		s3client, err = minio.NewWithCredentials(*s3endpoint, credential, true, *s3region)
		if err != nil {
			panic(err)
		}
		awsCredential := awscredentials.NewSharedCredentials(authfile, *awsProfile)
		signer = v4.NewSigner(awsCredential)
		signer.DisableHeaderHoisting = true
	} else {
		trashPath = filepath.Join(*rootDir, "/.Trash")
		if _, err := os.Stat(trashPath); err != nil {
			os.Mkdir(trashPath, 0744)
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

		// 		if req.URL.Path != KFS_CRT && req.URL.Path != KFS_KEY {
		w.Header().Add("Connection", "Keep-Alive")
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		if req.Method == "GET" {
			f := &fileHandler{Dir(*rootDir)}
			f.ServeHTTP(w, req)
		}
		// 		}
	})

	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	mux.Handle("/", handler)
	log.Println("Listen:", *listen, "Root:", *rootDir)
	// 	if enableTLS {
	// 		log.Println("Listen TLS:", *listenTLS, "Root:", *rootDir)
	// 	}
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
		/*
			if enableTLS {
				err = h2quic.ListenAndServeQUIC(*listen, crtPath, keyPath, LogHandler(gzipHandler(mux)))
				if err != nil {
					log.Println(err.Error())
					os.Exit(1)
				}
			}
		*/
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
