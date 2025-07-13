package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/cache"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
)

// GetListingsHandler handles GET /listings and returns all listings
type GetListingsHandler struct {
	Store *listings.Store
	Cache *cache.Cache
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
// @Header 		 200 {string} X-Cache  "HIT or MISS"
func (h GetListingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	cacheKey := "listings:all"
	if q != "" {
		cacheKey = "listings:q:" + q
	}

	var list []db.Listing
	if ok, _ := h.Cache.Get(r.Context(), cacheKey, &list); ok {
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(list)
		return
	}

	list, err := h.Store.ListByQuery(r.Context(), q)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	_ = h.Cache.Set(r.Context(), cacheKey, list)
	w.Header().Set("X-Cache", "MISS")
	_ = json.NewEncoder(w).Encode(list)
}
