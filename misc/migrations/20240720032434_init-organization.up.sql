CREATE TABLE IF NOT EXISTS "organizations" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "slug" varchar(255) UNIQUE NOT NULL,
    "description" TEXT NOT NULL,
    "avatar" varchar NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS "organization_member" (
    "pkid" bigserial PRIMARY KEY,
    "organization_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT NOT NULL,
    "role" varchar(50) NOT NULL CHECK (role IN ('admin', 'member', 'guest')),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    -- Foreign key constraints
    CONSTRAINT fk_organization
        FOREIGN KEY (organization_pkid) 
        REFERENCES "organizations" (pkid),

    CONSTRAINT fk_user
        FOREIGN KEY (user_pkid) 
        REFERENCES "users" (pkid)
);