CREATE TABLE IF NOT EXISTS "page_access_logs" (
    "pkid" BIGSERIAL PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT,
    "action" VARCHAR(25) NOT NULL,
    "last_accessed" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    
    CONSTRAINT "check_access_action" CHECK ("action" IN ('open', 'edit', 'upload')),
    
    CONSTRAINT fk_page_roles_page
        FOREIGN KEY ("page_pkid") 
        REFERENCES "pages" ("pkid") ON DELETE CASCADE,
    
    CONSTRAINT fk_page_roles_user
        FOREIGN KEY ("user_pkid")
        REFERENCES "users" ("pkid") ON DELETE CASCADE
);
