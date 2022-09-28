CREATE TABLE IF NOT EXISTS pizzas (
	id bigserial PRIMARY KEY,
	name text  NOT NULL DEFAULT 'untitled',
	reviewId int FOREIGN KEY reviewId REFERENCES reviews(id) ON UPDATE CASCADE ON DELETE CASCADE;
);