-- Modify "recipes" table
ALTER TABLE "recipes" ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" SET DEFAULT now();
-- Modify "relation_tuples" table
ALTER TABLE "relation_tuples" ALTER COLUMN "created_at" SET DEFAULT now();
-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "created_at" SET DEFAULT now(), ALTER COLUMN "updated_at" SET DEFAULT now();
