-- Drop FK constraint before type change
ALTER TABLE "recipes" DROP CONSTRAINT "recipes_user_id_fkey";
-- Modify "recipes" table
ALTER TABLE "recipes" ALTER COLUMN "user_id" TYPE character varying(26);
-- Modify "users" table
ALTER TABLE "users" ALTER COLUMN "id" TYPE character varying(26);
-- Restore FK constraint
ALTER TABLE "recipes" ADD CONSTRAINT "recipes_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
