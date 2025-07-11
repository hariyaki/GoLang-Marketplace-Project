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

// Create inserts a new Listing and handles created_at scanning
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
        // Try parsing RFC3339Nano, then fallback to common SQL datetime format
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

    var all []db.Listing
    for rows.Next() {
        var l db.Listing
        var rawCreated interface{}
        if err := rows.Scan(
            &l.ID,
            &l.Title,
            &l.Description,
            &l.PriceJPY,
            &rawCreated,
        ); err != nil {
            return nil, fmt.Errorf("scan listing: %w", err)
        }

        switch v := rawCreated.(type) {
        case time.Time:
            l.CreatedAt = v
        case string:
            if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
                l.CreatedAt = t
            } else if t, err := time.Parse("2006-01-02 15:04:05", v); err == nil {
                l.CreatedAt = t
            } else {
                return nil, fmt.Errorf("cannot parse created_at %q: %w", v, err)
            }
        default:
            return nil, fmt.Errorf("unexpected type %T for created_at", v)
        }

        all = append(all, l)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("rows error: %w", err)
    }

    return all, nil
}