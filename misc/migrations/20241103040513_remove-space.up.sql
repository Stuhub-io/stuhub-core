ALTER TABLE IF EXISTS "pages"
ADD COLUMN IF NOT EXISTS "org_pkid" BIGINT;

ALTER TABLE IF EXISTS "pages"
    DROP CONSTRAINT IF EXISTS "fk_page_organization";


ALTER TABLE IF EXISTS "pages"
    ADD CONSTRAINT "fk_page_organization"
    FOREIGN KEY (org_pkid)
    REFERENCES "organizations" (pkid)
    ON DELETE CASCADE;

-- Update the 'org_pkid' column in 'pages' table based on 'space_pkid'
DO $$ 
DECLARE column_exists BOOLEAN;
BEGIN
    -- Check if 'space_pkid' column exists in 'pages' table
    SELECT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'pages' 
        AND column_name = 'space_pkid'
    ) INTO column_exists;

    -- If column exists, run the update
    IF column_exists THEN
        UPDATE "pages" 
        SET "org_pkid" = spaces.org_pkid
        FROM spaces
        WHERE spaces.pkid = pages.space_pkid;
    END IF;
END $$;


ALTER TABLE IF EXISTS "pages"
DROP COLUMN IF EXISTS "space_pkid";

DROP TABLE IF EXISTS "space_member";
DROP TABLE IF EXISTS "spaces";