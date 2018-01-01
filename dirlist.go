/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : dirlist.go

* Purpose :

* Creation Date : 08-23-2017

* Last Modified : Mon 01 Jan 2018 12:08:21 AM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"encoding/json"
	"fmt"
	"github.com/coocood/freecache"
	"github.com/dustin/go-humanize"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

var (
	extIcon = map[string]string{
		"mp4":  "file-video-o",
		"mov":  "file-video-o",
		"wmv":  "file-video-o",
		"avi":  "file-video-o",
		"flv":  "file-video-o",
		"go":   "file-code-o",
		"mp3":  "file-audio-o",
		"jpeg": "file-image-o",
		"jpg":  "file-image-o",
		"png":  "file-image-o",
		"gif":  "file-image-o",
	}
	dirSizeCacheSize = 100 * 1024 * 1024
	dirSizeCache     = freecache.NewCache(dirSizeCacheSize)
)

const (
	staticTemplate = `
<html lang="en" ng-app="listApp">
<head>
<meta charset="UTF-8">
<meta name="referrer" content="none">
<meta name="google" content="notranslate">
<meta http-equiv="Content-Language" content="en">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/font-awesome/4.7.0/css/font-awesome.min.css">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/css/bootstrap.min.css" integrity="sha384-rwoIResjU2yc3z8GV/NPeZWAv56rSmLldC3R/AZzGRnGxQQKnKkoFVhFQhNUwEyJ" crossorigin="anonymous">
<script src="https://code.jquery.com/jquery-3.1.1.slim.min.js" integrity="sha384-A7FZj7v+d/sdmMqp/nOQwliLvUsJfDHW+k9Omg/a/EheAdgtzNs3hpfag6Ed950n" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/tether/1.4.0/js/tether.min.js" integrity="sha384-DztdAPBWPRXSA/3eYEEUWrWCy7G5KFbe8fFjk5JAIxUYHKkDx6Qin1DkWx51bBrb" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0-alpha.6/js/bootstrap.min.js" integrity="sha384-vBWWzlZJ8ea9aCX4pEW3rVHjgjt7zpkNpZk+02D9phzyeVkE+jo0ieGizqPLForn" crossorigin="anonymous"></script>
</head>
<style>
  body {
    font-family:"Microsoft Yahei","Helvetica Neue","Luxi Sans","DejaVu Sans",Tahoma,"Hiragino Sans GB",STHeiti;
  }
  table {
    font-size: 1.8em;
  }
  @media (max-width: 980px) {
    table {
      font-size: 1.8em;
    }
  }
</style>
<body ng-controller="listCtrl">

<div class="container">
  <div class="row">
    <div class="col-1">
	  <a href="[[.Url|urlBack|string]]#[[(printf "%s/" .Title)|hash]]"><h1> &lt; </h1></a>
    </div>
    <div class="col-5">
      <h1>[[.Title]]</h1>
    </div>
    <div class="col-2">
      <a href=[[urlSetQuery .Url "photo" "1"|string]]><button type="button" class="btn btn-secondary">PhotoGen</button></a>
    </div>
    <div class="col-4">
      <form action="[[.Url|string]]" method="get" class="bd-search hidden-sm-down">
	    <input type="text" class="form-control" placeholder="Search..." name="key" value="[[.Key]]" autofocus>
      </form>
    </div>
  </div>
</div>

<div class="container">
  <div class="row">
    <div class="col-11">
      <table class="table table-hover">
        <tr>
          <th><a href="[[index .Urls "name"]]">Name</a></th>
          <th><a href="[[index .Urls "size"]]">Size</a></th> 
          <th>Tags</th>
          <th>Func</th>
          <th><a href="[[index .Urls "lastMod"]]">LastMod</a></th>
        </tr>
        [[range .Files]]<tr[[if (index $.Meta.MetaInfo .Name).Label]] class="alert alert-[[(index $.Meta.MetaInfo .Name).Label]]"[[end]]>
		<td>[[if (index $.Meta.MetaInfo .Name).Star]]<i class="fa fa-star" aria-hidden="true"></i>  [[end]][[.Name|icon]]  <a name="[[.Name|hash]]" href="[[.Url|string]]">[[.Name]]</a></td>
          <td>[[.Size|size]]</td>
		  <td>[[if (index $.Meta.MetaInfo .Name).Tags]][[range (index $.Meta.MetaInfo .Name).Tags]][[.]] [[end]][[end]]</td>
          <td>
		  <div class="input-group">
		    <input type="checkbox" ng-model="enabled['[[.Name|hash]]']"><input ng-if="enabled['[[.Name|hash]]']" name="input" type="text" class="form-control" ng-model="file['[[.Name|hash]]']">
            <div class="dropdown" ng-if="enabled['[[.Name|hash]]']">
              <button class="btn btn-secondary dropdown-toggle" type="button" id="[[.Name|hash]]" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
                ...
              </button>
              <div class="dropdown-menu" aria-labelledby="[[.Name|hash]]">
            
			  <a style="color:red;" class="dropdown-item" href="[[urlSetQuery $.Url "delete" "1" "name" .Name|string]]">Delete</a>
            
            	<a class="dropdown-item" ng-if="file['[[.Name|hash]]']" href="[[urlSetQuery $.Url "rename" "1" "name" .Name|string]]&rename={{file['[[.Name|hash]]']}}">Rename to {{file['[[.Name|hash]]']}}</a>
            
            	<a class="dropdown-item" ng-if="file['[[.Name|hash]]']" href="[[urlSetQuery $.Url "rename" "1" "name" .Name|string]]&rename={{file['[[.Name|hash]]']}}[[.Name]]">Rename to {{file['[[.Name|hash]]']}}[[.Name]]</a>
            
            	<a class="dropdown-item" ng-if="file['[[.Name|hash]]']" href="[[urlSetQuery $.Url "rename" "1" "name" .Name|string]]&name=[[.Name]]{{file['[[.Name|hash]]']}}">Rename to [[.Name]]{{file['[[.Name|hash]]']}}</a>
            
            	<a class="dropdown-item" ng-if="file['[[.Name|hash]]']" href="[[urlSetQuery $.Url "addtags" "1" "name" .Name|string]]&tags={{file['[[.Name|hash]]']}}">add tags {{file['[[.Name|hash]]']}}</a>
            	<a class="dropdown-item" ng-if="file['[[.Name|hash]]']" href="[[urlSetQuery $.Url "updatetags" "1" "name" .Name|string]]&tags={{file['[[.Name|hash]]']}}">update tags {{file['[[.Name|hash]]']}}</a>
            	<a class="dropdown-item" href="[[urlSetQuery $.Url "uncompress" "1" "name" .Name|string]]&pass={{file['[[.Name|hash]]']}}">uncompress with pass {{file['[[.Name|hash]]']}}</a>
            
            	<a class="dropdown-item" href="[[urlSetQuery $.Url "star" "1" "name" .Name|string]]"><i class="fa fa-star" aria-hidden="true"></i></a>

            	<a class="dropdown-item" href="[[urlSetQuery $.Url "setlabel" "0" "name" .Name|string]]">label 0</a>
            	<a class="dropdown-item alert alert-success" href="[[urlSetQuery $.Url "setlabel" "success" "name" .Name|string]]">label 1</a>
            	<a class="dropdown-item alert alert-info" href="[[urlSetQuery $.Url "setlabel" "info" "name" .Name|string]]">label 2</a>
            	<a class="dropdown-item alert alert-warning" href="[[urlSetQuery $.Url "setlabel" "warning" "name" .Name|string]]">label 3</a>
            	<a class="dropdown-item alert alert-danger" href="[[urlSetQuery $.Url "setlabel" "danger" "name" .Name|string]]">label 4</a>

              </div>
            </div>
		  </div>
		  </td>
          <td>[[.ModTime|time]]</td>
        </tr>[[end]]
      </table>
    </div>
    <div class="col-1">
    </div>
  </div>
</div>

<div class="container">
  <div class="row">
    <div class="col-1">
	  <a href="[[.Url|urlBack|string]]#[[(printf "%s/" .Title)|hash]]"><h1> &lt; </h1></a>
	</div>
    <div class="col-11">
	</div>
  </div>
</div>

<script type="text/javascript" src="https://ajax.googleapis.com/ajax/libs/angularjs/1.6.5/angular.min.js"></script>
<script>
[[.NgScript]]
</script>
</body>
</html>
`
	ngScript = `
angular.module("listApp", [])
.controller("listCtrl", function($scope, $http, $location, $window, $interval) {
	$scope.file = {};
	$scope.enabled = {};
})
`
)

type Page struct {
	Title    string
	NgScript template.JS
	Files    []*PageFile
	Url      url.URL
	Urls     map[string]string
	Key      string
	Desc     string
	Meta     Meta
}

type PageFile struct {
	Name    string
	Url     url.URL
	Size    uint64
	LastMod string
	ModTime time.Time
}

type Meta struct {
	Root     string
	MetaInfo map[string]MetaInfo
}

func (m *Meta) Get() error {
	if m.MetaInfo == nil {
		m.MetaInfo = make(map[string]MetaInfo)
	}
	metaFile := filepath.Join(m.Root, ".KFS_META")
	b, err := ioutil.ReadFile(metaFile)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	err = json.Unmarshal(b, m)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return nil
}
func (m *Meta) Write() error {
	metaFile := filepath.Join(m.Root, ".KFS_META")
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Println(err.Error())
		return err
	}
	return ioutil.WriteFile(metaFile, b, 0644)
}

type MetaInfo struct {
	Label  string
	Tags   []string
	Star   bool
	OldLoc string
}

var tmpFundMap = template.FuncMap{
	"string":        tmpString,
	"urlSetQuery":   urlSetQuery,
	"urlCleanQuery": urlCleanQuery,
	"urlBack":       urlBack,
	"time":          prettyTime,
	"icon":          getIcon,
	"hash":          hash,
	"size": func(i uint64) template.HTML {
		return template.HTML(humanize.IBytes(i))
	},
}

func tmpString(i interface{}) template.HTML {
	switch v := i.(type) {
	case url.URL:
		return template.HTML(v.String())
	case *url.URL:
		return template.HTML(v.String())
	}
	return "-"
}

func urlBack(u url.URL) url.URL {
	// 	return url.URL{Path: filepath.Dir(u.Path)}
	p := strings.Split(u.Path, "/")
	var b string
	if len(p) > 2 {
		b = "/" + strings.Join(p[1:len(p)-2], "/") + "/"
		if strings.Contains(b, "//") {
			b = "/"
		}
	}
	u.Path = b
	return u
}

func urlCleanQuery(u url.URL) url.URL {
	u.RawQuery = ""
	return u
}

// func urlSetQuery(u url.URL, key, value string) url.URL {
func urlSetQuery(u url.URL, kv ...string) url.URL {
	if len(kv)%2 == 1 {
		return u
	}
	v := u.Query()
	for i := 0; i < len(kv); i += 2 {
		v.Set(kv[i], kv[i+1])
	}
	u.RawQuery = v.Encode()
	return u
}

func prettyTime(t time.Time) template.HTML {
	since := time.Since(t)
	switch {
	case since < (1 * time.Second):
		return template.HTML("1s")

	case since < (60 * time.Second):
		s := strings.Split(fmt.Sprint(since), ".")[0]
		return template.HTML(s + "s")

	case since < (60 * time.Minute):
		s := strings.Split(fmt.Sprint(since), ".")[0]
		return template.HTML(strings.Split(s, "m")[0] + "m")

	case since < (24 * time.Hour):
		s := strings.Split(fmt.Sprint(since), ".")[0]
		return template.HTML(strings.Split(s, "h")[0] + "h")

	default:
		return template.HTML(t.Format("01-02-06"))
	}
	return template.HTML("")
}

func dirList1(w http.ResponseWriter, f http.File, r *http.Request, filedir string) {
	dirs, err := f.Readdir(-1)
	if err != nil {
		// TODO: log err.Error() to the Server.ErrorLog, once it's possible
		// for a handler to get at its Server via the ResponseWriter. See
		// Issue 12438.
		http.Error(w, "Error reading directory", http.StatusInternalServerError)
		return
	}

	v := r.URL.Query()
	doPhoto := v.Get("photo")
	if doPhoto == "1" {
		v.Del("photo")
		r.URL.RawQuery = v.Encode()
		renderPhoto(w, r, filedir)
		return
	}

	meta := Meta{Root: filedir}
	var createMeta bool
	err = meta.Get()
	if err != nil {
		createMeta = true
	}

	doSetLabel := v.Get("setlabel")
	if len(doSetLabel) != 0 {
		name := v.Get("name")
		m, ok := meta.MetaInfo[name]
		if !ok {
			m = MetaInfo{}
		}
		if doSetLabel != "0" {
			m.Label = doSetLabel
		} else {
			m.Label = ""
		}
		meta.MetaInfo[name] = m
		meta.Write()
		v.Del("setlabel")
		v.Del("name")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = hash(name)
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}

	doStar := v.Get("star")
	if len(doStar) != 0 {
		name := v.Get("name")
		m, ok := meta.MetaInfo[name]
		if !ok {
			m = MetaInfo{}
		}
		if m.Star {
			m.Star = false
		} else {
			m.Star = true
		}
		meta.MetaInfo[name] = m
		meta.Write()
		v.Del("star")
		v.Del("name")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = hash(name)
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}
	doDelete := v.Get("delete")
	if len(doDelete) != 0 {
		oldname := v.Get("name")
		newname := oldname
		_, err := os.Stat(filepath.Join(trash, newname))
		var i int
		for err == nil {
			n := oldname
			if strings.HasSuffix(n, "/") {
				n = n[:len(n)-1]
			}
			newname = fmt.Sprintf("%s_%d", n, i)
			_, err = os.Stat(filepath.Join(trash, newname))
			i++
		}
		if err != nil {
			log.Println("do mv", filepath.Join(filedir, oldname), filepath.Join(trash, newname))
			os.Rename(filepath.Join(filedir, oldname), filepath.Join(trash, newname))
		}
		m, ok := meta.MetaInfo[oldname]
		if ok {
			delete(meta.MetaInfo, oldname)
		}
		meta.Write()

		m2 := Meta{Root: trash}
		m2.Get()
		m.OldLoc = filepath.Join(filedir, oldname)
		m2.MetaInfo[newname] = m
		m2.Write()

		v.Del("delete")
		v.Del("name")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = ""
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}
	doAddTags := v.Get("addtags")
	if len(doAddTags) != 0 {
		name := v.Get("name")
		m, ok := meta.MetaInfo[name]
		if !ok {
			m = MetaInfo{}
		}
		for _, v := range strings.Split(v.Get("tags"), " ") {
			v = strings.Trim(v, " ")
			v = strings.ToLower(v)
			add := true
			for _, exist := range m.Tags {
				if v == exist {
					add = false
				}
			}
			if add {
				m.Tags = append(m.Tags, v)
			}
		}
		sort.Strings(m.Tags)
		meta.MetaInfo[name] = m
		meta.Write()
		v.Del("addtags")
		v.Del("tags")
		v.Del("name")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = hash(name)
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}
	doUpdateTags := v.Get("updatetags")
	if len(doUpdateTags) != 0 {
		name := v.Get("name")
		m, ok := meta.MetaInfo[name]
		if !ok {
			m = MetaInfo{}
		}
		m.Tags = []string{}
		if v.Get("tags") != "-" {
			for _, v := range strings.Split(v.Get("tags"), " ") {
				v = strings.Trim(v, " ")
				v = strings.ToLower(v)
				add := true
				for _, exist := range m.Tags {
					if v == exist {
						add = false
					}
				}
				if add {
					m.Tags = append(m.Tags, v)
				}
			}
		}
		sort.Strings(m.Tags)
		meta.MetaInfo[name] = m
		meta.Write()
		v.Del("updatetags")
		v.Del("tags")
		v.Del("name")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = hash(name)
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}

	doUncompress := v.Get("uncompress")
	if len(doUncompress) != 0 {
		name := v.Get("name")
		pass := v.Get("pass")
		uncompress(filedir, name, pass)
		v.Del("uncompress")
		v.Del("name")
		v.Del("pass")
		r.URL.RawQuery = v.Encode()
		r.URL.Fragment = hash(name)
		http.Redirect(w, r, r.URL.String(), 302)
		return
	}

	orderBy := v.Get("by")
	desc := v.Get("desc")
	key := v.Get("key")
	var list []os.FileInfo
	if len(key) != 0 {
		for _, v := range dirs {
			if v.Name() == key {
				u := url.URL{Path: v.Name()}
				http.Redirect(w, r, u.String(), 302)
				return
			}
			if strings.Contains(v.Name(), key) {
				list = append(list, v)
			}
		}
		dirs = list
	}

	r.URL.RawQuery = v.Encode()

	page := new(Page)
	stat, _ := f.Stat()
	page.Title = stat.Name()
	if r.URL.Path == "/" {
		page.Title = "/"
	}
	page.NgScript = ngScript
	page.Url = *r.URL
	page.Urls = make(map[string]string)
	page.Key = key

	for _, t := range []string{"name", "size", "lastMod"} {
		v.Set("by", t)
		switch desc {
		case "1":
			v.Set("desc", "0")
		default:
			v.Set("desc", "1")
		}
		page.Urls[t] = "?" + v.Encode()
	}

	switch orderBy {
	case "name":
		sort.Slice(dirs, func(i, j int) bool {
			if desc == "0" {
				return dirs[i].Name() < dirs[j].Name()
			}
			return dirs[i].Name() > dirs[j].Name()
		})
	case "size":
		sort.Slice(dirs, func(i, j int) bool {
			if desc == "0" {
				return dirSize(filepath.Join(filedir, dirs[i].Name())) < dirSize(filepath.Join(filedir, dirs[j].Name()))
			}
			return dirSize(filepath.Join(filedir, dirs[i].Name())) > dirSize(filepath.Join(filedir, dirs[j].Name()))
		})
	default:
		sort.Slice(dirs, func(i, j int) bool {
			if desc == "0" {
				return dirs[i].ModTime().Unix() < dirs[j].ModTime().Unix()
			}
			return dirs[i].ModTime().Unix() > dirs[j].ModTime().Unix()
		})
	}

	//
	// 	metaFile := filepath.Join(filedir, ".KFS_META")
	page.Meta = meta
	// 	meta := Meta{MetaInfo: make(map[string]MetaInfo)}
	// 	b, err := ioutil.ReadFile(metaFile)
	// 	json.Unmarshal(b, meta.MetaInfo)
	//
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, d := range dirs {
		if d.Name() == ".KFS_META" {
			continue
		}
		var f PageFile
		f.Name = d.Name()
		if d.IsDir() {
			f.Name += "/"
			f.Size = dirSize(filepath.Join(filedir, d.Name()))
		} else {
			f.Size = uint64(d.Size())
		}
		u := url.URL{Path: f.Name}
		f.Url = u
		f.ModTime = d.ModTime()

		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		// 		url := url.URL{Path: name}
		page.Files = append(page.Files, &f)
	}
	if createMeta {
		meta.Root = filedir
		for _, v := range page.Files {
			meta.MetaInfo[v.Name] = MetaInfo{}
		}
		b, _ := json.MarshalIndent(meta, "", "  ")
		ioutil.WriteFile(filepath.Join(meta.Root, ".KFS_META"), b, 0644)
	}
	tmpl, err := template.New("page").Funcs(tmpFundMap).Delims(`[[`, `]]`).Parse(staticTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = tmpl.Execute(w, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getIcon(i interface{}) template.HTML {
	switch v := i.(type) {
	case string:
		if strings.HasSuffix(v, "/") {
			return template.HTML(`<i class="fa fa-folder-open-o" aria-hidden="true"></i>`)
		}
		ext := filepath.Ext(v)
		if v, ok := extIcon[ext]; ok {
			return template.HTML(fmt.Sprintf(`<i class="fa fa-%s" aria-hidden="true"></i>`, v))
		}
	}
	return `<i class="fa fa-file-o" aria-hidden="true"></i>`
}

func dirSize(path string) uint64 {
	// 	t1 := time.Now()
	if b, err := dirSizeCache.Get([]byte(path)); err == nil {
		// 		log.Println("size HIT", path, string(b))
		u, _ := strconv.ParseUint(string(b), 10, 64)
		return u
	}
	var size uint64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if info != nil && !info.IsDir() {
			size += uint64(info.Size())
		}
		return err
	})
	if err != nil {
		log.Println(err.Error())
	}
	dirSizeCache.Set([]byte(path), []byte(fmt.Sprint(size)), 60*30)
	// 	log.Println("size MISS", path, s, time.Since(t1))
	return size
}
