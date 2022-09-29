CREATE TABLE IF NOT EXISTS pizzas (
	id bigserial PRIMARY KEY,
	name text  NOT NULL DEFAULT 'untitled',
	review_id int,
	CONSTRAINT review_fk 
		FOREIGN KEY (review_id) 
			REFERENCES reviews(id) 
			ON UPDATE CASCADE 
			ON DELETE CASCADE
);