package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
)

type PostListingHandler struct {
	Store *listings.Store
}

type postListingRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	PriceJPY    int64  `json:"price_jpy"`
}

// CreateListing godoc
// @Summary      Create a new listing
// @Description  Inserts a listing row and returns the created object.
// @Tags         listings
// @Accept       json
// @Produce      json
// @Param        payload  body      postListingRequest  true  "Listing payload"
// @Success      201      {object}  db.Listing
// @Failure      400      {string}  string  "invalid JSON or missing fields"
// @Failure      500      {string}  string  "database error"
// @Router       /listings [post]
func (h PostListingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var req postListingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}
	if req.Title == "" || req.Description == "" || req.PriceJPY < 0 {
		http.Error(w, "missing or invalid fields", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	l, err := h.Store.Create(ctx, db.Listing{
		Title:       req.Title,
		Description: req.Description,
		PriceJPY:    req.PriceJPY,
	})
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(l)
}
