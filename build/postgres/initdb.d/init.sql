CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
  "pkid" bigint PRIMARY KEY,
  "id" uuid UNIQUE NOT NULL DEFAULT (uuid_generate_v4()),
  "username" varchar UNIQUE NOT NULL,
  "password" varchar,
  "activate_email" varchar UNIQUE,
  "gmail" varchar UNIQUE,
  "first_name" varchar,
  "last_name" varchar,
  "avatar" varchar,
  "created_at" timestamp DEFAULT (now()),
  "updated_at" timestamp DEFAULT (now()),
  "is_activated" bool
);

CREATE UNIQUE INDEX ON "users" ("pkid");