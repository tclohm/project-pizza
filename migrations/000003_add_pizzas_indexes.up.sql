CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE EXTENSION IF NOT EXISTS btree_gin;
CREATE INDEX IF NOT EXISTS pizzas_name_idx ON pizzas USING GIN (to_tsvector('simple', name));