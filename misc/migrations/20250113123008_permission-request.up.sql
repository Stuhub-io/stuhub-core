CREATE TABLE IF NOT EXISTS "page_permission_request_log" (
    "pkid" bigserial PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT,
    "email" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "status" VARCHAR(20) DEFAULT 'pending' NOT NULL CHECK ("status" IN ('pending', 'approved', 'rejected')),

    CONSTRAINT fk_page_permission_request_log_page
        FOREIGN KEY (page_pkid) 
        REFERENCES "pages" (pkid) ON DELETE CASCADE,
    
    CONSTRAINT fk_page_permission_request_log_user
        FOREIGN KEY (user_pkid)
        REFERENCES "users" (pkid) ON DELETE CASCADE
)