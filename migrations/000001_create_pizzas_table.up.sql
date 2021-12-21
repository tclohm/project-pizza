CREATE TABLE IF NOT EXISTS pizzas (
	id bigserial PRIMARY KEY,
	name text  NOT NULL DEFAULT 'untitled',
	style text NOT NULL DEFAULT 'unknown',
	price double precision NOT NULL, 
	description text NOT NULL DEFAULT '...',
	cheesiness double precision NOT NUll DEFAULT 0,
	flavor double precision NOT NULL DEFAULT 0,
	sauciness double precision NOT NULL DEFAULT 0,
	saltiness double precision NOT NULL DEFAULT 0,
	charness double precision NOT NULL DEFAULT 0,
	created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
);