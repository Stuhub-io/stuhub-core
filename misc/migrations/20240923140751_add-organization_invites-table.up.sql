CREATE TABLE IF NOT EXISTS "organization_invites" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() NOT NULL,
    "organization_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "expired_at" TIMESTAMP WITH TIME ZONE NOT NULL,

    CONSTRAINT fk_owner
        FOREIGN KEY (user_pkid) 
        REFERENCES "users" (pkid) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_invite_pkid ON "organization_invites" (id);