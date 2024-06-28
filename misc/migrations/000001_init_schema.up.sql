CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "users" (
  "pkid" bigserial PRIMARY KEY,
  "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
  
  "email" varchar(255) NOT NULL,
  "password" varchar(128),

  "first_name" varchar(255) NOT NULL,
  "last_name" varchar(255) NOT NULL,
  "avatar"  varchar NOT NULL,
  
  "oath_gmail" varchar NOT NULL,

  "activated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
  "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT (now())
);