CREATE TABLE IF NOT EXISTS images {
	id bigserial PRIMARY KEY,
	filename text NOT NULL,
	type text NOT NULL,
	location text NOT NULL,
	created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
}