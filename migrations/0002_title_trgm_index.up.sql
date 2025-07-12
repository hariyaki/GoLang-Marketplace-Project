CREATE INDEX listings_title_idx ON listings USING gin (title gin_trgm_ops);
