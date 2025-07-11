// internal/handlers/get_listings_test.go
package handlers_test

import (
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

func TestGetListingsHandler_OK(t *testing.T) {
	// --- Arrange ---
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	store := listings.NewStore(dbMock)

	// Prepare fake rows
	rows := sqlmock.NewRows([]string{
		"id", "title", "description", "price_jpy", "created_at",
	}).
		AddRow(1, "foo", "first item", int64(1000), "2025-07-10T12:00:00Z").
		AddRow(2, "bar", "second item", int64(2000), "2025-07-11T13:30:00Z")

	// Expect the SELECT query
	mock.
		ExpectQuery(regexp.QuoteMeta(
			`SELECT id, title, description, price_jpy, created_at
    FROM listings
    ORDER BY created_at DESC;`,
		)).
		WillReturnRows(rows)

	handler := handlers.GetListingsHandler{Store: store}

	req := httptest.NewRequest(http.MethodGet, "/listings", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusOK, rec.Code, "expected status 200 OK")

	var got []db.Listing
	require.NoError(t, json.NewDecoder(rec.Body).Decode(&got), "should decode JSON response")
	require.Len(t, got, 2, "should return two listings")

	//Verify first listing
	require.Equal(t, int64(1), got[0].ID)
	require.Equal(t, "foo", got[0].Title)
	require.Equal(t, "first item", got[0].Description)
	require.Equal(t, int64(1000), got[0].PriceJPY)

	//Verify second listing
	require.Equal(t, int64(2), got[1].ID)
	require.Equal(t, "bar", got[1].Title)
	require.Equal(t, "second item", got[1].Description)
	require.Equal(t, int64(2000), got[1].PriceJPY)

	//Ensure all expectations were met
	require.NoError(t, mock.ExpectationsWereMet())
}
