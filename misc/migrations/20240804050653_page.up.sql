
CREATE TABLE IF NOT EXISTS "pages" (
    "pkid" bigserial PRIMARY KEY,
    "id" UUID DEFAULT uuid_generate_v4() UNIQUE NOT NULL,
    "name" varchar(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "space_pkid" BIGINT NOT NULL,
    "parent_page_pkid" BIGINT,
    
    "view_type" varchar(50) NOT NULL CHECK (view_type IN ('document', 'table')),
    
    CONSTRAINT fk_page_space
        FOREIGN KEY (space_pkid) 
        REFERENCES "spaces" (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_parent
        FOREIGN KEY (parent_page_pkid)
        REFERENCES "pages" (pkid)
        ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_parent_page_pkid ON "pages" (parent_page_pkid);
