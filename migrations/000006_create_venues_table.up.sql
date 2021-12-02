CREATE TABLE IF NOT EXISTS venues (
	id bigserial PRIMARY KEY,
	name text NOT NULL,
	lat double precision NOT NULL,
	lon double precision NOT NULL,
	address text NOT NULL
);