-- Create enum type "visibility"
CREATE TYPE "visibility" AS ENUM ('public', 'private', 'restricted');
-- Modify "recipes" table
ALTER TABLE "recipes" ADD COLUMN "visibility" "visibility" NOT NULL;
-- Create "relation_tuples" table
CREATE TABLE "relation_tuples" (
  "id" uuid NOT NULL,
  "object_type" character varying(255) NOT NULL,
  "object_id" character varying(255) NOT NULL,
  "relation" character varying(255) NOT NULL,
  "user_id" character varying(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "unique" UNIQUE ("object_type", "object_id", "relation", "user_id")
);
