CREATE TABLE IF NOT EXISTS "activity" (
    pkid BIGINT PRIMARY KEY,
    actor_pkid BIGINT,
    page_pkid BIGINT,
    action_code TEXT,
    label TEXT,
    metadata TEXT,  -- JSONB is not supported, store JSON as a string
    created_at TIMESTAMP
);