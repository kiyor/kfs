/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : log.go

* Purpose :

* Creation Date : 08-27-2017

* Last Modified : Sun Aug 27 17:44:33 2017

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
)

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	w.length = len(b)
	return w.ResponseWriter.Write(b)
}

func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t1 := time.Now()
		ctx := context.Background()
		writer := statusWriter{w, 0, 0}

		r = r.WithContext(ctx)
		next.ServeHTTP(&writer, r)

		rang := r.Header.Get("Range")
		if len(rang) == 0 {
			rang = "-"
		}

		res := fmt.Sprintf("%v %v %v %v %v %v %v", r.RemoteAddr, writer.status, writer.length, r.Method, r.Host+r.RequestURI, rang, time.Since(t1))
		log.Println(res)

	})
}
