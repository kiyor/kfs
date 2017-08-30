/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : photo.go

* Purpose :

* Creation Date : 08-28-2017

* Last Modified : Tue 29 Aug 2017 11:57:31 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
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

<link rel="shortcut icon" href="/favicon.ico" type="image/x-icon" />
<script type="text/javascript" src="https://dev.2ns.io/wnacg/jquery-3.1.0.min.js"></script>
<script type="text/javascript" src="https://dev.2ns.io/wnacg/scroll.photos.js"></script>
<script type="text/javascript">var hash = location.hash;
if(!hash){
	hash = 0;
}else{
	hash = parseInt(hash.replace("#","")) - 1;
}
var imglist = [{{.|imageslist}}];
$(function(){
	imgscroll.beLoad($("#img_list"),imglist,hash)
});</script>
</head><body style="background:#444;">
	<div id="page_scale"></div>
	<div id="cite_vote" style="background-color:#fafafa;display:none;">
		<div class="toolbar" style="display:block;">
			<a href="javascript:;" onclick="citeShow(0);" class="back"></a>
			<div class="title"></div>
		</div>
		<div id="sns_cite_vote_list" style="margin-top:40px;"></div>
	</div>
	<div class="mask_panel" id="mask_panel" style="position:fixed;">
	</div>
		<div class="shareBox" id="shareBox">
		</div>	
	</div>
	<div id="img_list">
		<style>
		.adBox{text-align:center;padding:0 3px 25px;}
		.adBox img{max-width:100% !important;height:auto !important;}
		</style>
	</div>
	<div id="img_load" style="text-align:center;color:#AAA;"><img src="https://dev.2ns.io/wnacg/loading.gif" /><br /><span>少女讀取中...</span></div>
	<div class="section4" id="control_block" style="display:none;">
	</div>
</body>
</html>`

type Image string
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

func (images Images) List() string {
	var out string
	for _, v := range images.Images {
		v = Image(strings.Replace(string(v), " ", "%20", -1))
		out += fmt.Sprintf("{url:\"%s\"},", v)
	}
	if len(out) > 0 {
		return out[:len(out)-1]
	}
	return ""
}

func imagesList(images Images) string {
	return images.List()
}

func init() {
	flag.Parse()
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
