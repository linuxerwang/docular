package server

import (
	"archive/zip"
	"io"
	"net/http"
	"path/filepath"
)

func (wh *wrapperHandler) serveMaff(w http.ResponseWriter, r *http.Request, maffPath, filePath string) {
	path := filepath.Join(string(wh.docroot), maffPath)

	// MAFF file is actually a zip file. Read from it.

	reader, err := zip.OpenReader(path)
	if err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer reader.Close()

	maffDir := ""
	for _, f := range reader.File {
		// All archived files are put in a random folder, and the folder is the first
		// archived element.

		if f.FileInfo().IsDir() && maffDir == "" {
			maffDir = f.Name
			continue
		}

		if f.Name == filepath.Join(maffDir, filePath) {
			rc, err := f.Open()
			if err != nil {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			w.WriteHeader(http.StatusOK)
			io.Copy(w, rc)
			rc.Close()

			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}
