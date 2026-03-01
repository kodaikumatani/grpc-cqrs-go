-- Create "recipes" table
CREATE TABLE "recipes" (
  "id" uuid NOT NULL,
  "user_id" character varying(26) NOT NULL,
  "title" character varying(255) NOT NULL,
  "description" character varying(255) NOT NULL,
  "created_at" timestamp NOT NULL,
  "updated_at" timestamp NOT NULL,
  PRIMARY KEY ("id")
);
