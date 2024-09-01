CREATE TABLE IF NOT EXISTS "documents" (
    "pkid" bigserial PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL UNIQUE,
    "content" TEXT NOT NULL,
    "json_content" JSON,
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),

    CONSTRAINT fk_document_page
        FOREIGN KEY (page_pkid)
        REFERENCES "pages" (pkid)
        ON DELETE CASCADE
);

DO $$
DECLARE
    page_pkid BIGINT;  -- Declare a scalar variable for pkid
BEGIN
    -- Loop through each page where view_type is 'document'
    FOR page_pkid IN
        SELECT pkid FROM pages WHERE view_type = 'document'
    LOOP
        -- Insert corresponding document entry
        INSERT INTO documents (page_pkid, content, json_content)
        VALUES (page_pkid, '', NULL);
    END LOOP;
END $$;
