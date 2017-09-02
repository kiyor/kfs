/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : photo.go

* Purpose :

* Creation Date : 08-28-2017

* Last Modified : Sat 02 Sep 2017 07:12:37 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	// 	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

var goodExt = []string{".jpg", ".png", ".gif", ".jpeg"}

var index = `<!DOCTYPE html>
<html lang="en-US">
<head>
<meta name="viewport" content="width=device-width,initial-scale=1.0,maximum-scale=2.5,user-scalable=yes"/>  
<meta charset="UTF-8"/>
<meta name="full-screen" content="yes"/>
<meta name="browsermode" content="application"/>
<meta name="x5-pagetype" content="webapp"/>
<meta name="format-detection" content="telephone=no"/>
<meta name="apple-mobile-web-app-capable" content="yes"/>
<meta name="apple-mobile-web-app-status-bar-style" content="white"/>
<title>{{.Title}}</title>
<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon"/>
<style>
#topBtn {
  display: none;
  position: fixed;
  bottom: 20px;
  right: 30px;
  z-index: 99;
  border: none;
  outline: none;
  background-color: #111;
  color: white;
  cursor: pointer;
  padding: 15px;
  border-radius: 10px;
}
#topBtn:hover {
  background-color: #111;
}
</style>
</head>
<body style="background:#444;">
	<div id="img_list">
	</div>
    <button onclick="topFunction()" id="topBtn" title="Go to top">Top</button>
	<div id="img_load" style="text-align:center;color:#AAA;"><img src="https://dev.2ns.io/wnacg/loading.gif" /><br /><span>少女讀取中...</span></div>

	<script type="text/javascript" src="https://dev.2ns.io/wnacg/jquery-3.1.0.min.js"></script>
	<script type="text/javascript" src="https://dev.2ns.io/wnacg/scroll.photos.js"></script>
	<script type="text/javascript">
	var hash = location.hash;
	if(!hash){
		hash = 0;
	}else{
		hash = parseInt(hash.replace("#","")) - 1;
	}
	var imglist = [{{.|imageslist}}];
	$(function(){
		imgscroll.beLoad($("#img_list"),imglist,hash)
	});
	</script>
	<script>
	// When the user scrolls down 20px from the top of the document, show the button
	window.onscroll = function() {scrollFunction()};
	
	function scrollFunction() {
	    if (document.body.scrollTop > 20 || document.documentElement.scrollTop > 20) {
	        document.getElementById("topBtn").style.display = "block";
	    } else {
	        document.getElementById("topBtn").style.display = "none";
	    }
	}
	
	// When the user clicks on the button, scroll to the top of the document
	function topFunction() {
	    document.body.scrollTop = 0;
	    document.documentElement.scrollTop = 0;
	}
	</script>
</body>
</html>`

type Image template.HTML

type Images struct {
	Title  string
	Images []Image
}

func (images Images) Exec() string {
	f := template.FuncMap{
		"imageslist": imagesList,
	}
	t, err := template.New("index").Funcs(f).Parse(index)
	if err != nil {
		panic(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, images)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func (images Images) List() template.JS {
	var out string
	for _, v := range images.Images {
		// 		v = Image(strings.Replace(string(v), " ", "%20", -1))
		s := (&url.URL{Path: string(v)}).String()
		out += fmt.Sprintf("{url:\"%s\"},", s)
	}
	if len(out) > 0 {
		return template.JS(out[:len(out)-1])
	}
	return ""
}

func imagesList(images Images) template.JS {
	return images.List()
}

func init() {
	flag.Parse()
}

func renderPhoto(w http.ResponseWriter, r *http.Request, dir string) {
	var images Images
	abs, _ := filepath.Abs(dir)
	images.Title = filepath.Base(abs)
	fs := readDir(dir)
	if len(fs) == 0 {
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}
	for _, f := range fs {
		images.Images = append(images.Images, Image(f.Name()))
	}
	f := template.FuncMap{
		"imageslist": imagesList,
	}
	t, err := template.New("index").Funcs(f).Parse(index)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
	err = t.Execute(w, images)
	if err != nil {
		fmt.Fprintln(w, err.Error())
	}
}

func mkphoto(dir string) error {
	var images Images
	abs, _ := filepath.Abs(dir)
	images.Title = filepath.Base(abs)
	fs := readDir(dir)
	if len(fs) == 0 {
		return nil
	}
	for _, f := range fs {
		images.Images = append(images.Images, Image(f.Name()))
	}

	return ioutil.WriteFile(dir+"/photo.html", []byte(images.Exec()), 0644)
}

func readDir(path string) (fs []os.FileInfo) {
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		for _, ext := range goodExt {
			if strings.ToLower(filepath.Ext(f.Name())) == ext {
				fs = append(fs, f)
				break
			}
		}
	}
	return
}
