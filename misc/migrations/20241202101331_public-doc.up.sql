


CREATE TABLE IF NOT EXISTS "public_token" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    "page_pkid" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "archived_at" TIMESTAMP WITH TIME ZONE,

    CONSTRAINT fk_public_token_page
        FOREIGN KEY (page_pkid)
        REFERENCES "pages" (pkid) ON DELETE CASCADE
)
