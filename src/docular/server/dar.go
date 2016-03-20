package server

import (
	"io"
	"net/http"
	"path/filepath"

	"docular/dar"
)

func (wh *wrapperHandler) serveDar(w http.ResponseWriter, r *http.Request, darPath, filePath string) {
	path := filepath.Join(string(wh.docroot), darPath)

	// Read from the dar file

	reader, err := dar.OpenReader(path)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer reader.Close()

	for _, f := range reader.File {
		if f.Name == filePath {
			rc, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.WriteHeader(http.StatusOK)
			io.Copy(w, rc)
			rc.Close()
		}
	}
}
