CREATE TABLE IF NOT EXISTS pizzas (
	id bigserial PRIMARY KEY,
	name text  NOT NULL,
	style text NOT NULL,
	description text NOT NULL,
	cheesiness double precision NOT NUll,
	flavor double precision NOT NULL,
	sauciness double precision NOT NULL,
	saltiness double precision NOT NULL,
	charness double precision NOT NULL,
	image_filename text NOT NULL,
	image_content_type text NOT NULL,
	image_location text NOT NULL,
	created_at TIMESTAMP(0) with time zone NOT NULL DEFAULT NOW()
);