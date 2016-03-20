package server

import (
	"net/http"
	"path"
	"strings"
)

func (wh *wrapperHandler) serveStaticFile(w http.ResponseWriter, r *http.Request, fpath string) {
	if fpath == "/_/default.css" {
		http.ServeFile(w, r, path.Join(Webstatic, "styles/default.css"))
		return
	}

	if fpath == "/_/jquery.min.js" {
		http.ServeFile(w, r, path.Join(Webstatic, "scripts/jquery.min.js"))
		return
	}

	if fpath == "/_/default.js" {
		http.ServeFile(w, r, path.Join(Webstatic, "scripts/default.js"))
		return
	}

	if strings.HasPrefix(fpath, "/_/fonts/") {
		http.ServeFile(w, r, path.Join(Webstatic, "fonts", fpath[9:]))
		return
	}
}
