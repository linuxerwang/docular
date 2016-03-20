package server

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	htmlHeader = `<head><meta charset="utf-8">
<link href="/_/default.css" rel="stylesheet">
<script type="text/javascript" src="/_/jquery.min.js"></script>
<script type="text/javascript" src="/_/default.js"></script>
</head>
`
	htmlPwd  = `<div class="docular-pwd">Current Position: %s</div>`
	htmlRoot = `<div class="docular-entry docular-dir docular-parent" data-url="%s">↰</div>`
	htmlDir  = `<div class="docular-entry docular-dir" data-url="%s">%s</div>`
	htmlFile = `<div class="docular-entry docular-file" data-url="%s">%s</div>`

	darFile  = `<div class="docular-entry docular-dar" data-url="%s/">%s</div>`
	maffFile = `<div class="docular-entry docular-maff" data-url="%s/">%s</div>`
)

var (
	Webstatic   string
	filesToShow map[string]struct{}
)

type ByName []os.FileInfo

func (bn ByName) Len() int {
	return len(bn)
}

func (bn ByName) Swap(i, j int) {
	bn[i], bn[j] = bn[j], bn[i]
}

func (bn ByName) Less(i, j int) bool {
	return bn[i].Name() < bn[j].Name()
}

type wrapperHandler struct {
	docroot http.Dir
	wrapped http.Handler
}

func (wh *wrapperHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Maybe serve a file in MAFF archive.

	pos := strings.Index(r.URL.Path, ".maff/")
	if pos > 0 {
		maffPath := r.URL.Path[:pos+5]
		filePath := r.URL.Path[pos+5:]

		if strings.HasSuffix(filePath, "/") {
			filePath += "index.html"
		}

		wh.serveMaff(w, r, maffPath, filePath)
		return
	}

	// Maybe serve a file in dar archive.

	pos = strings.Index(r.URL.Path, ".dar/")
	if pos > 0 {
		darPath := r.URL.Path[:pos+4]
		filePath := r.URL.Path[pos+4:]

		if strings.HasSuffix(filePath, "/") {
			filePath += "index.html"
		}

		wh.serveDar(w, r, darPath, filePath)
		return
	}

	// If url ends with "/" (requesting directory), serve it.
	if strings.HasSuffix(r.URL.Path, "/") {
		wh.serveDir(w, r)
		return
	}

	if strings.HasSuffix(r.URL.Path, "/index.html") {
		// Handle index.html specifically, otherwise net/http would redirect to .../

		fpath := filepath.Clean(filepath.Join(string(wh.docroot), r.URL.Path))
		f, err := os.Open(fpath)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		defer f.Close()

		d, err1 := f.Stat()
		if err1 != nil {
			http.NotFound(w, r)
			return
		}

		if d.IsDir() {
			if checkLastModified(w, r, d.ModTime()) {
				return
			}
			wh.serveDir(w, r)
			return
		}

		http.ServeContent(w, r, d.Name(), d.ModTime(), f)

		return
	}

	if strings.HasPrefix(r.URL.Path, "/_/") {
		wh.serveStaticFile(w, r, r.URL.Path)
		return
	}

	if strings.HasSuffix(r.URL.Path, ".md") {
		wh.serveMarkDown(w, r, r.URL.Path)
		return
	}

	// Ignore requests for hidden files.

	base := filepath.Base(r.URL.Path)
	if strings.HasPrefix(base, ".") {
		http.NotFound(w, r)
		return
	}

	// Otherwise delegate to the wrapped handler.
	wh.wrapped.ServeHTTP(w, r)
}

// modtime is the modification time of the resource to be served, or IsZero().
// return value is whether this request is now complete.
func checkLastModified(w http.ResponseWriter, r *http.Request, modtime time.Time) bool {
	if modtime.IsZero() {
		return false
	}

	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	if t, err := time.Parse(http.TimeFormat, r.Header.Get("If-Modified-Since")); err == nil && modtime.Before(t.Add(1*time.Second)) {
		h := w.Header()
		delete(h, "Content-Type")
		delete(h, "Content-Length")
		w.WriteHeader(http.StatusNotModified)
		return true
	}
	w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	return false
}

func (wh *wrapperHandler) serveDir(w http.ResponseWriter, r *http.Request) {
	showParentDir := r.URL.Path != "/"
	showFullPage := r.Method == "GET"

	dir, err := wh.docroot.Open(r.URL.Path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dir.Close()

	items, err := dir.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		return
	}

	if showFullPage {
		w.Write([]byte(htmlHeader))
		w.Write([]byte(`<body>`))
	}

	w.Write([]byte(fmt.Sprintf(htmlPwd, r.URL.Path)))

	folders, files, dars, maffs := groupFiles(items)

	// Show directories

	if len(folders) > 0 || showParentDir {
		w.Write([]byte("<div class=\"docular-container\">"))
		if showParentDir {
			p := filepath.Dir(r.URL.Path[:len(r.URL.Path)-1])
			if !strings.HasSuffix(p, "/") {
				p += "/"
			}
			w.Write([]byte(fmt.Sprintf(htmlRoot, p)))
		}

		for _, f := range folders {
			w.Write([]byte(fmt.Sprintf(htmlDir, r.URL.Path+f.Name()+"/", f.Name())))
		}

		w.Write([]byte("</div>"))
	}

	// Show html/markdown files

	if len(files) > 0 {
		w.Write([]byte("<div class=\"docular-container\">"))
		for _, f := range files {
			w.Write([]byte(fmt.Sprintf(htmlFile, r.URL.Path+f.Name(), f.Name())))
		}
		w.Write([]byte("</div>"))
	}

	// Show dar files

	if len(dars) > 0 {
		w.Write([]byte("<div class=\"docular-container\">"))
		for _, f := range dars {
			w.Write([]byte(fmt.Sprintf(darFile, r.URL.Path+f.Name(), f.Name())))
		}
		w.Write([]byte("</div>"))
	}

	// Show maff files

	if len(maffs) > 0 {
		w.Write([]byte("<div class=\"docular-container\">"))
		for _, f := range maffs {
			w.Write([]byte(fmt.Sprintf(maffFile, r.URL.Path+f.Name(), f.Name())))
		}
		w.Write([]byte("</div>"))
	}
	if showFullPage {
		w.Write([]byte(`</body>`))
	}
}

func groupFiles(items []os.FileInfo) ([]os.FileInfo, []os.FileInfo, []os.FileInfo, []os.FileInfo) {
	folders := []os.FileInfo{}
	files := []os.FileInfo{}
	dars := []os.FileInfo{}
	maffs := []os.FileInfo{}

	for _, f := range items {
		if f.IsDir() {
			if !strings.HasPrefix(f.Name(), ".") {
				folders = append(folders, f)
			}
		} else {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			if ext == ".html" || ext == ".htm" || ext == ".md" {
				files = append(files, f)
			} else if ext == ".dar" {
				dars = append(dars, f)
			} else if ext == ".maff" {
				maffs = append(maffs, f)
			}
		}
	}

	sort.Sort(ByName(folders))
	sort.Sort(ByName(files))
	sort.Sort(ByName(dars))
	sort.Sort(ByName(maffs))

	return folders, files, dars, maffs
}

func NewWrapperHandler(docroot http.Dir) *wrapperHandler {
	return &wrapperHandler{
		docroot: docroot,
		wrapped: http.FileServer(http.Dir(docroot)),
	}
}
