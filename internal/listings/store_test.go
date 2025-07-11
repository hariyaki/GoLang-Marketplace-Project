package listings_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
	"github.com/hariyaki/GoLang-Marketplace-Project/internal/listings"
	"github.com/stretchr/testify/require"
)

func TestStore_Create(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	store := listings.NewStore(dbMock)
	want := db.Listing{
		Title:       "test",
		Description: "desc",
		PriceJPY:    123,
	}

	mock.ExpectQuery(regexp.QuoteMeta(
		"INSERT INTO listings (title, description, price_jpy) VALUES ($1, $2, $3) RETURNING id, created_at;",
	)).
		WithArgs(want.Title, want.Description, want.PriceJPY).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).AddRow(1, "2025-07-11T00:00:00Z"))

	got, err := store.Create(context.Background(), want)
	require.NoError(t, err)
	require.Equal(t, int64(1), got.ID)
	require.Equal(t, want.Title, got.Title)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestStore_List(t *testing.T) {
	dbMock, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer dbMock.Close()

	store := listings.NewStore(dbMock)

	rows := sqlmock.NewRows([]string{"id", "title", "description", "price_jpy", "created_at"}).
		AddRow(1, "t1", "d1", 100, "2025-07-10T00:00:00Z").
		AddRow(2, "t2", "d2", 200, "2025-07-11T00:00:00Z")

	mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT id, title, description, price_jpy, created_at FROM listings ORDER BY created_at DESC;",
	)).WillReturnRows(rows)

	got, err := store.List(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 2)
	require.Equal(t, int64(1), got[0].ID)
	require.Equal(t, int64(2), got[1].ID)
	require.NoError(t, mock.ExpectationsWereMet())
}
