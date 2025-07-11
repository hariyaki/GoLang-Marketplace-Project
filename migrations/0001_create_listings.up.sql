CREATE TABLE listings (
    id          BIGSERIAL PRIMARY KEY,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL,
    price_jpy   BIGINT      NOT NULL CHECK (price_jpy >= 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
