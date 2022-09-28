ALTER TABLE reviews ADD CONSTRAINT reviews_cheesiness_check CHECK (cheesiness BETWEEN 0 AND 5);

ALTER TABLE reviews ADD CONSTRAINT reviews_flavor_check CHECK (flavor BETWEEN 0 AND 5);

ALTER TABLE reviews ADD CONSTRAINT reviews_sauciness_check CHECK (sauciness BETWEEN 0 AND 5);

ALTER TABLE reviews ADD CONSTRAINT reviews_saltiness_check CHECK (saltiness BETWEEN 0 AND 5);

ALTER TABLE reviews ADD CONSTRAINT reviews_charness_check CHECK (charness BETWEEN 0 AND 5);

ALTER TABLE reviews ADD CONSTRAINT reviews_spiciness_check CHECK (spiciness BETWEEN 0 AND 5);