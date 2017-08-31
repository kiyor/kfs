/* -.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.-.

* File Name : dirlist.go

* Purpose :

* Creation Date : 08-23-2017

* Last Modified : Thu 31 Aug 2017 11:12:08 PM UTC

* Created By : Kiyor

_._._._._._._._._._._._._._._._._._._._._.*/

package main

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
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
)

const (
	staticTemplate = `
<html lang="en">
<head>
<meta charset="UTF-8">
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
<body>

<div class="container">
  <div class="row">
    <div class="col-1">
      <a href="{{.Url|urlBack|string}}#{{(printf "%s/" .Title)|hash}}"><h1> &lt; </h1></a>
    </div>
    <div class="col-5">
      <h1>{{.Title}}</h1>
    </div>
    <div class="col-2">
      <a href={{urlSetQuery .Url "photo" "1"|string}}><button type="button" class="btn btn-secondary">PhotoGen</button></a>
    </div>
    <div class="col-4">
      <form action="{{.Url|string}}" method="get" class="bd-search hidden-sm-down">
        <input type="text" name="key" placeholder="Search..." value="{{.Key}}" autofocus>
      </form>
    </div>
  </div>
</div>

<div class="container">
  <div class="row">
    <div class="col-11">
      <table class="table table-hover">
        <tr>
          <th><a href="{{index .Urls "name"}}">Name</a></th>
          <th><a href="{{index .Urls "size"}}">Size</a></th> 
          <th>Func</th>
          <th><a href="{{index .Urls "lastMod"}}">LastMod</a></th>
        </tr>
        {{range .Files}}<tr>
          <td>{{.Name|icon}}  <a name="{{.Name|hash}}" href="{{.Url|string}}">{{.Name}}</a></td>
          <td>{{.Size}}</td>
          <td>
<div class="dropdown">
  <button class="btn btn-secondary dropdown-toggle" type="button" id="{{.Name|hash}}" data-toggle="dropdown" aria-haspopup="true" aria-expanded="false">
    Func
  </button>
  <div class="dropdown-menu" aria-labelledby="{{.Name|hash}}">
    <a class="dropdown-item" href="{{urlSetQuery .Url "delete" "1"|string}}">Delete</a>
  </div>
</div>
		  </td>
          <td>{{.ModTime|time}}</td>
        </tr>{{end}}
      </table>
    </div>
    <div class="col-1">
    </div>
  </div>
</div>

<div class="container">
  <div class="row">
    <div class="col-1">
      <a href="{{.Url|urlBack|string}}#{{(printf "%s/" .Title)|hash}}"><h1> &lt; </h1></a>
	</div>
    <div class="col-11">
	</div>
  </div>
</div>


</body>
</html>
`
)

type Page struct {
	Title string
	Files []*PageFile
	Url   url.URL
	Urls  map[string]string
	Key   string
	Desc  string
}

type PageFile struct {
	Name    string
	Url     url.URL
	Size    string
	LastMod string
	ModTime time.Time
}

var tmpFundMap = template.FuncMap{
	"string":        tmpString,
	"urlSetQuery":   urlSetQuery,
	"urlCleanQuery": urlCleanQuery,
	"urlBack":       urlBack,
	"time":          prettyTime,
	"icon":          getIcon,
	"hash":          hash,
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

func urlSetQuery(u url.URL, key, value string) url.URL {
	v := u.Query()
	v.Set(key, value)
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

func dirList1(w http.ResponseWriter, f http.File, r *http.Request, dir string) {
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
		mkphoto(dir)
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
				// 			b, err := filepath.Match(v.Name(), key)
				// 			if err != nil {
				// 				log.Println(err.Error())
				// 			}
				// 			if b {
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
	page.Url = *r.URL
	page.Urls = make(map[string]string)
	page.Key = key
	// 	v.Set("photo", "1")
	// 	r.URL.RawQuery = v.Encode()
	// 	page.UrlPhoto = r.URL.String()

	for _, t := range []string{"name", "size", "lastMod"} {
		v.Set("by", t)
		switch desc {
		case "1":
			v.Set("desc", "0")
		default:
			v.Set("desc", "1")
		}
		// 		r.URL.RawQuery = v.Encode()
		// 		page.Urls[t] = r.url.string()
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
				return dirs[i].Size() < dirs[j].Size()
			}
			return dirs[i].Size() > dirs[j].Size()
		})
	default:
		sort.Slice(dirs, func(i, j int) bool {
			if desc == "0" {
				return dirs[i].ModTime().Unix() < dirs[j].ModTime().Unix()
			}
			return dirs[i].ModTime().Unix() > dirs[j].ModTime().Unix()
		})
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	for _, d := range dirs {
		var f PageFile
		f.Name = d.Name()
		if d.IsDir() {
			f.Name += "/"
			// 			f.Icon = `<i class="fa fa-folder-open-o" aria-hidden="true"></i>`
		} else {
			// 			f.Icon = getIcon(f.Name)
		}
		// 		f.Name = htmlReplacer.Replace(f.Name)
		u := url.URL{Path: f.Name}
		f.Url = u
		f.Size = humanize.IBytes(uint64(d.Size()))
		// 		f.LastMod = d.ModTime().Format("01-02-2006 15:04:05")
		f.ModTime = d.ModTime()

		// name may contain '?' or '#', which must be escaped to remain
		// part of the URL path, and not indicate the start of a query
		// string or fragment.
		// 		url := url.URL{Path: name}
		page.Files = append(page.Files, &f)
	}
	tmpl, err := template.New("page").Funcs(tmpFundMap).Parse(staticTemplate)
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
