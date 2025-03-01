CREATE TABLE IF NOT EXISTS page_star (
    "pkid" BIGSERIAL PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL,
    "user_pkid" BIGINT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "order" FLOAT DEFAULT 0 NOT NULL,

    CONSTRAINT fk_page_star_user
        FOREIGN KEY (user_pkid)
        REFERENCES users (pkid) ON DELETE CASCADE,

    CONSTRAINT fk_page_star_page
        FOREIGN KEY (page_pkid)
        REFERENCES pages (pkid) ON DELETE CASCADE
);

CREATE UNIQUE INDEX unique_user_page_idx ON page_star (user_pkid, page_pkid);