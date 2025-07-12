package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
)

type GetListingHandler struct {
	Store *listings.Store
}

func (h GetListingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/listings/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	l, err := h.Store.GetByID(r.Context(), id)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(l)
}
