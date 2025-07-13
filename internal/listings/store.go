// internal/listings/store.go
package listings

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
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

type ListOpts struct {
	Query	string
	Limit	int
	Offset	int
	Sort	string
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
func (s *Store) List(ctx context.Context, o ListOpts) ([]db.Listing, error) {
	var (
		args	[]any
		sb		strings.Builder
	)
	sb.WriteString(`
		SELECT id,title,description,price_jpy,image_url,created_at
        FROM listings
	`)

	if o.Query != "" {
		args = append(args, o.Query)
		sb.WriteString(` WHERE title ILIKE '%'||$1||'%'`)
	}

	switch o.Sort {
	case "price_asc":
		sb.WriteString(` ORDER BY price_jpy ASC`)
	case "price_desc":
		sb.WriteString(` ORDER BY price_jpy DESC`)
	default:
		sb.WriteString(` ORDER BY created_at DESC`)
	}

	limPos := len(args) + 1
	offPos := limPos + 1
	args = append(args, o.Limit, o.Offset)
	sb.WriteString(fmt.Sprintf(` LIMIT $%d OFFSET $%d`, limPos, offPos))

	rows, err := s.DB.QueryContext(ctx, sb.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []db.Listing
	for rows.Next() {
		var l db.Listing

		if err := rows.Scan(&l.ID, &l.Title, &l.Description, &l.PriceJPY, &l.ImageURL, &l.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, l)
	}

	return out, rows.Err()
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

func (s *Store) UpdateImage(ctx context.Context, id int64, url string) (db.Listing, error) {
	const q = `
		UPDATE listings SET image_url=$2
		WHERE id=$1
		RETURNING id,title,description,price_jpy,image_url,created_at`
	var l db.Listing
	err := s.DB.QueryRowContext(ctx, q, id, url).
		Scan(&l.ID, &l.Title, &l.Description, &l.PriceJPY, &l.ImageURL, &l.CreatedAt)
	return l, err
}

func IsNotFound(err error) bool {
	return errors.Is(err, sql.ErrNoRows)
}
