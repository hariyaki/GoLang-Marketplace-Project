package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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
// @Summary      List or search listings
// @Tags         listings
// @Produce      json
// @Param        q       query   string false  "search keyword"
// @Param        limit   query   int    false  "max results (1-100)" minimum(1) maximum(100)
// @Param        offset  query   int    false  "offset for pagination"
// @Param        sort    query   string false  "new|price_asc|price_desc"
// @Success      200 {array} db.Listing
// @Failure      500 {string} string
// @Header       200 {string} X-Cache "HIT or MISS"
// @Router       /listings [get]
func (h GetListingsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	lim := parseInt(r.URL.Query().Get("limit"), 20)
	off := parseInt(r.URL.Query().Get("offset"), 0)
	sort := r.URL.Query().Get("sort")

	if lim <= 0 || lim > 100 {
		lim = 20
	}
	if off < 0 {
		off = 0
	}

	cacheKey := fmt.Sprintf("listings:q%s:lim=%d:off=%d:sort=%s", q, lim, off, sort)
	var list []db.Listing
	if ok, _ := h.Cache.Get(r.Context(), cacheKey, &list); ok {
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(list)
		return
	}

	opts := listings.ListOpts{
		Query:  q,
		Limit:  lim,
		Offset: off,
		Sort:   sort,
	}

	list, err := h.Store.List(r.Context(), opts)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	_ = h.Cache.Set(r.Context(), cacheKey, list)
	w.Header().Set("X-Cache", "MISS")
	_ = json.NewEncoder(w).Encode(list)
}

func parseInt(s string, def int) int {
	if v, err := strconv.Atoi(s); err == nil {
		return v
	}
	return def
}
