package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/shurcooL/github_flavored_markdown"
)

func (wh *wrapperHandler) serveMarkDown(w http.ResponseWriter, r *http.Request, path string) {
	path = filepath.Join(string(wh.docroot), r.URL.Path)

	f, err := os.Open(path)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Convert to HTML
	unsafe := github_flavored_markdown.Markdown(b)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	printPageHeader(w)

	_, err = w.Write([]byte(html))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	printPageFooter(w)

}

func printPageHeader(w http.ResponseWriter) {
	w.Write([]byte(`<html>
<header>
<meta charset="utf-8">
<link href="/_/default.css" rel="stylesheet">
</header><body><div class="doc-view">`))
}

func printPageFooter(w http.ResponseWriter) {
	w.Write([]byte("</div></body></html>"))
}
