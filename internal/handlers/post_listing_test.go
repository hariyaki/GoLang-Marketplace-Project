package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/handlers"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (http.Handler, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)

	store := listings.NewStore(dbMock)
	h := handlers.PostListingHandler{Store: store}
	return h, mock
}

func TestPostListingHandler_OK(t *testing.T) {
	h, mock := setup(t)

	body := []byte(`{"title":"shoes", "description":"running shoes", "price_jpy":6000}`)
	req := httptest.NewRequest(http.MethodPost, "/listings", bytes.NewReader(body))
	rec := httptest.NewRecorder()

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO listings (title, description, price_jpy) VALUES ($1, $2, $3) RETURNING id, created_at;",
	)).
		WithArgs("shoes", "running shoes", int64(6000)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(42, "2025-07-11T00:00:00Z"))

	h.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)

	var got db.Listing
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got))
	require.Equal(t, int64(42), got.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
