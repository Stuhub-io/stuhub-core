CREATE TABLE IF NOT EXISTS "users" (
  "id" bigserial PRIMARY KEY,
  "username" varchar NOT NULL,
  "email" varchar NOT NULL,
  "password" varchar,
  "avatar"  varchar,
  "is_oauth" boolean DEFAULT false,
  "is_activated" boolean DEFAULT false,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now())
);