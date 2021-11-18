ALTER TABLE pizzas ADD CONSTRAINT pizzas_cheesiness_check CHECK (cheesiness BETWEEN 0 AND 5);

ALTER TABLE pizzas ADD CONSTRAINT pizzas_flavor_check CHECK (flavor BETWEEN 0 AND 5);

ALTER TABLE pizzas ADD CONSTRAINT pizzas_sauciness_check CHECK (sauciness BETWEEN 0 AND 5);

ALTER TABLE pizzas ADD CONSTRAINT pizzas_saltiness_check CHECK (saltiness BETWEEN 0 AND 5);

ALTER TABLE pizzas ADD CONSTRAINT pizzas_charness_check CHECK (charness BETWEEN 0 AND 5);