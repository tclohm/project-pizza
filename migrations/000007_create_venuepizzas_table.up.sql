CREATE TABLE IF NOT EXISTS venuepizzas (
	id bigserial PRIMARY KEY,
	venue_id int,
	pizza_id int,
	CONSTRAINT venue_fk
	 FOREIGN KEY (venue_id) REFERENCES venues(id) ON UPDATE CASCADE ON DELETE CASCADE,
	CONSTRAINT pizza_fk
	 FOREIGN KEY (pizza_id) REFERENCES pizzas(id) ON UPDATE CASCADE ON DELETE CASCADE
);