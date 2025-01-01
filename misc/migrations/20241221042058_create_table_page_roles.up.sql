
CREATE TABLE IF NOT EXISTS "page_roles" (
    "pkid" bigserial PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT,
    "email" VARCHAR(255) NOT NULL,
    "role" VARCHAR(20) DEFAULT 'viewer' NOT NULL  CHECK ("role" IN ('viewer', 'editor')),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    
    CONSTRAINT fk_page_roles_page
        FOREIGN KEY (page_pkid) 
        REFERENCES "pages" (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_page_roles_user
        FOREIGN KEY (user_pkid)
        REFERENCES "users" (pkid) ON DELETE CASCADE
);
