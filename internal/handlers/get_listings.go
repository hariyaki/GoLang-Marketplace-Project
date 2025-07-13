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

// ListOrSearchListings godoc
// @Summary      List all listings or search by title
// @Description  When the optional **q** query parameter is supplied, the results are filtered (case-insensitive substring match on *title*).
// @Tags         listings
// @Produce      json
// @Param        q   query     string  false  "Title search keyword"
// @Success      200  {array}   db.Listing
// @Failure      500  {string}  string  "database error"
// @Router       /listings [get]
func (h GetListingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	list, err := h.Store.ListByQuery(r.Context(), q)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
