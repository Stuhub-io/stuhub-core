CREATE TABLE IF NOT EXISTS "spaces" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    "org_pkid" BIGINT NOT NULL,
    "name" varchar(255) NOT NULL,
    "description" TEXT NOT NULL,
    "is_private" BOOLEAN NOT NULL DEFAULT FALSE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_space_org
        FOREIGN KEY (org_pkid) 
        REFERENCES "organizations" (pkid) ON DELETE CASCADE
);  -- Removed the trailing comma here

CREATE TABLE IF NOT EXISTS "space_member" (
    "pkid" bigserial PRIMARY KEY,
    "space_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT NOT NULL,
    "role" varchar(50) NOT NULL CHECK (role IN ('owner', 'member', 'guest')),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- Foreign key constraints
    CONSTRAINT fk_space
        FOREIGN KEY (space_pkid) 
        REFERENCES "spaces" (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_user
        FOREIGN KEY (user_pkid) 
        REFERENCES "users" (pkid) ON DELETE CASCADE
);