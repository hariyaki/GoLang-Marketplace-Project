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
// Updated to handle dynamic types for created_at in both Create and List

// Store wraps *sql.DB for mock tests
var _ = fmt.Errorf // ensure fmt is imported even if not used elsewhere

type Store struct {
	DB *sql.DB
}

func NewStore(db *sql.DB) *Store { return &Store{DB: db} }

func parseCreated(raw interface{}) (time.Time, error) {
	switch v := raw.(type) {
	case time.Time:
		return v, nil

	case string:
		return parseTimeString(v)

	case []byte: // ← sqlmock uses this path
		return parseTimeString(string(v))

	default:
		return time.Time{}, fmt.Errorf("unexpected type %T for created_at", v)
	}
}

// parseTimeString tries a couple of common layouts.
func parseTimeString(s string) (time.Time, error) {
	if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
		return t, nil
	}
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t, nil
	}
	return time.Time{}, fmt.Errorf("cannot parse timestamp %q", s)
}

// Create inserts a new Listing and handles created_at scanning
func (s *Store) Create(ctx context.Context, l db.Listing) (db.Listing, error) {
	const q = `
        INSERT INTO listings (title, description, price_jpy)
        VALUES ($1, $2, $3)
        RETURNING id, created_at;
    `

	// Scan into interface{} for handling both time.Time and string
	var rawCreated interface{}
	if err := s.DB.
		QueryRowContext(ctx, q, l.Title, l.Description, l.PriceJPY).
		Scan(&l.ID, &rawCreated); err != nil {
		return l, fmt.Errorf("insert listing: %w", err)
	}

	t, err := parseCreated(rawCreated)
	if err != nil {
		return l, err
	}
	l.CreatedAt = t
	return l, nil
}

// List returns all listings ordered by created_at desc and handles created_at scanning
func (s *Store) List(ctx context.Context) ([]db.Listing, error) {
	const q = `
        SELECT id, title, description, price_jpy, created_at
        FROM listings
        ORDER BY created_at DESC;
    `

	rows, err := s.DB.QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("query listings: %w", err)
	}
	defer rows.Close()

	var out []db.Listing
	for rows.Next() {
		var l db.Listing
		var rawCreated interface{}

		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.PriceJPY, &rawCreated); err != nil {
			return nil, fmt.Errorf("scan listing: %w", err)
		}

		t, err := parseCreated(rawCreated)
		if err != nil {
			return nil, err
		}
		l.CreatedAt = t
		out = append(out, l)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return out, nil
}

// GetByID returns either a single listing or sql.ErrNoRows
func (s *Store) GetByID(ctx context.Context, id int64) (db.Listing, error) {
	const q = `SELECT id,title,description,price_jpy,created_at FROM listings WHERE id=$1`
	var l db.Listing
	err := s.DB.QueryRowContext(ctx, q, id).
		Scan(&l.ID, &l.Title, &l.Description, &l.PriceJPY, &l.CreatedAt)
	return l, err
}

func (s *Store) ListByQuery(ctx context.Context, qstr string) ([]db.Listing, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if qstr == "" {
		rows, err = s.DB.QueryContext(ctx,
			`SELECT id,title,description,price_jpy,created_at
			FROM listings WHERE title ILIKE '%'||$1||'%' ORDER BY created_at DESC`,
			qstr)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []db.Listing
	for rows.Next() {
		var l db.Listing
		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.PriceJPY, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}
	return out, rows.Err()
}
