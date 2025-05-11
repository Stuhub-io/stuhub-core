
CREATE TABLE IF NOT EXISTS "activity" (
    "pkid" BIGSERIAL PRIMARY KEY,
    "user_pkid" BIGINT NOT NULL,
    "action_code" VARCHAR(100) NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "snapshot" JSONB DEFAULT '{}' NOT NULL,
    
    CONSTRAINT fk_activity_user
        FOREIGN KEY (user_pkid)
        REFERENCES users (pkid) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "relate_page_activity" (
    "pkid" BIGSERIAL PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL,
    "activity_pkid" BIGINT NOT NULL,

    CONSTRAINT fk_relate_page_activity_page
        FOREIGN KEY (page_pkid)
        REFERENCES pages (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_relate_page_activity_activity
        FOREIGN KEY (activity_pkid)
        REFERENCES activity (pkid) ON DELETE CASCADE
);

-- Uinque together 
CREATE UNIQUE INDEX IF NOT EXISTS "unique_relate_page_activity_idx"
    ON "relate_page_activity" (page_pkid, activity_pkid);