package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
)

// GetListingsHandler handles GET /listings and returns all listings
type GetListingsHandler struct {
	Store *listings.Store
}

func (h GetListingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	all, err := h.Store.List(ctx)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(all)
}
