CREATE TABLE IF NOT EXISTS "spaces" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    "owner_id" BIGINT NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" TEXT NOT NULL,
    "is_private" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- Foreign key constraints
    CONSTRAINT fk_owner
        FOREIGN KEY (owner_id) 
        REFERENCES "users" (pkid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "space_member" (
    "pkid" bigserial PRIMARY KEY,
    "space_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- Foreign key constraints
    CONSTRAINT fk_space
        FOREIGN KEY (space_pkid) 
        REFERENCES "spaces" (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_user
        FOREIGN KEY (user_pkid) 
        REFERENCES "users" (pkid)
);
