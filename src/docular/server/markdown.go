package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

const (
  renderHtmlFlags = 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_FRACTIONS |
		blackfriday.HTML_SMARTYPANTS_DASHES |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES |
		blackfriday.HTML_TOC

	renderExtensions = 0 |
		blackfriday.EXTENSION_AUTO_HEADER_IDS |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_FOOTNOTES |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_TITLEBLOCK
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
	unsafe := renderMarkdown(b)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	printPageHeader(w)

	_, err = w.Write([]byte(html))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}

	printPageFooter(w)
}

func renderMarkdown(input []byte) []byte {
  // set up the HTML renderer
	renderer := blackfriday.HtmlRenderer(renderHtmlFlags, "", "")
	return blackfriday.MarkdownOptions(input, renderer, blackfriday.Options{
		Extensions: renderExtensions})
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
