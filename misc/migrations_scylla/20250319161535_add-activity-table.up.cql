CREATE TABLE IF NOT EXISTS "activity" (
    actor_pkid BIGINT,
    page_pkid BIGINT,
    action_code TEXT,
    label TEXT,
    metadata TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((actor_pkid), created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);