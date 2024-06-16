CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "password" varchar NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now())
);