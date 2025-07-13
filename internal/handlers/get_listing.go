package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/cache"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
)

type GetListingHandler struct {
	Store *listings.Store
	Cache *cache.Cache
}

// GetListing godoc
// @Summary      Retrieve a single listing
// @Tags         listings
// @Produce      json
// @Param        id   path      int  true  "Listing ID"
// @Success      200  {object}  db.Listing
// @Failure      400  {string}  string  "invalid id"
// @Failure      404  {string}  string  "not found"
// @Failure      500  {string}  string  "database error"
// @Router       /listings/{id} [get]
// @Header 		 200 {string} X-Cache  "HIT or MISS"
func (h GetListingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/listings/")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	l, err := h.Store.GetByID(r.Context(), id)

	key := fmt.Sprintf("listing:%d", id)
	if ok, _ := h.Cache.Get(r.Context(), key, &l); ok {
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(l)
		return
	}

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	_ = h.Cache.Set(r.Context(), key, l)
	w.Header().Set("X-Cache", "MISS")
	_ = json.NewEncoder(w).Encode(l)
}
