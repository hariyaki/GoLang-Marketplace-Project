// internal/listings/store.go
package listings

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

	// Scan into interface{} for handling both time.Time and string
	var rawCreated interface{}
	row := s.DB.QueryRowContext(ctx, q, l.Title, l.Description, l.PriceJPY)
	if err := row.Scan(&l.ID, &rawCreated); err != nil {
		return l, fmt.Errorf("insert listing: %w", err)
	}

	// Convert rawCreated into a time.Time
	switch v := rawCreated.(type) {
	case time.Time:
		l.CreatedAt = v

	case string:
		// Try RFC3339Nano first, then fallback to common SQL datetime format
		if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
			l.CreatedAt = t
		} else if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
			l.CreatedAt = t
		} else {
			return l, fmt.Errorf("cannot parse created_at %q: %w", v, err)
		}

	default:
		return l, fmt.Errorf("unexpected type %T for created_at", v)
	}

	return l, nil
}
