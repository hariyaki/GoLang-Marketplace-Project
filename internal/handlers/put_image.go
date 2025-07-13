package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/storage"
)

type PutImageHandler struct {
	Store   *listings.Store
	Storage *storage.FS
}

func (h PutImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/listings/")
	idStr = strings.TrimSuffix(idStr, "/image")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
	}

	file, hdr, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := path.Ext(hdr.Filename)
	filename := fmt.Sprintf("%d%s", id, ext)
	url, err := h.Storage.Save(r.Context(), filename, file)
	if err != nil {
		http.Error(w, "store error", http.StatusInternalServerError)
		return
	}

	l, err := h.Store.UpdateImage(r.Context(), id, url)
	if err != nil {
		status := http.StatusInternalServerError
		if listings.IsNotFound(err) {
			status = http.StatusNotFound
		}
		http.Error(w, err.Error(), status)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(l)
}
