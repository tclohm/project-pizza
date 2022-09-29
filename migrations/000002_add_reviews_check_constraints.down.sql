ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_cheesiness_check;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_flavor_check;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_sauciness_check;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_saltiness_check;

ALTER TABLE reviews DROP CONSTRAINT IF EXISTS reviews_charness_check;