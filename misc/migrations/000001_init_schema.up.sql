CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "pkid" bigserial PRIMARY KEY,
  "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
  "username" varchar(255) NOT NULL,
  "email" varchar(255) NOT NULL,
  "password" varchar(128),
  "avatar"  varchar NOT NULL,
  "is_oauth" boolean DEFAULT false NOT NULL,
  "is_activated" boolean DEFAULT false NOT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now())
);