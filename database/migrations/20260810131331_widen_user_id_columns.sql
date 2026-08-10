-- Modify "recipes" table
ALTER TABLE "recipes" ALTER COLUMN "user_id" TYPE character varying(255);
-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "id" TYPE character varying(255);
