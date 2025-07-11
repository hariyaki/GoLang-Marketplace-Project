package listings

import (
	"context"
	"database/sql"

	"github.com/hariyaki/GoLang-Marketplace-Project/internal/db"
)

// Store wraps *sql.DB for mock tests
type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{DB: db} }

func (s *Store) Create(ctx context.Context, l db.Listing) (db.Listing, error) {
	const q = `
		INSERT INTO listings (title, description, price_jpy)
		VALUES ($1, $2, $3)
		RETURNING id, created_at;
	`

	err := s.DB.QueryRowContext(ctx, q, l.Title, l.Description, l.PriceJPY).
		Scan(&l.ID, &l.CreatedAt)

	return l, err
}
