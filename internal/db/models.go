package db

import "time"

type Listing struct {
	ID          int64     `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	PriceJPY    int64     `db:"price_jpy" json:"price_jpy"`
	ImageURL    *string   `db:"image_url"   json:"image_url,omitempty"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
}
