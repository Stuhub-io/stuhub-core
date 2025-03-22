
-- 1. Recreate the `spaces` table (if needed, adjust columns as per the original schema)
CREATE TABLE IF NOT EXISTS "spaces" (
    pkid BIGINT PRIMARY KEY,
    org_pkid BIGINT,
    -- Add other necessary columns for `spaces` here
    FOREIGN KEY (org_pkid) REFERENCES "organizations" (pkid) ON DELETE CASCADE
);

-- 2. Recreate the `space_member` table (adjust columns as needed)
CREATE TABLE IF NOT EXISTS "space_member" (
    -- Define the schema for `space_member` based on original setup
);

-- 3. Re-add the `space_pkid` column in `pages`
ALTER TABLE IF EXISTS "pages"
ADD COLUMN IF NOT EXISTS "space_pkid" BIGINT;

-- 4. Set `space_pkid` in `pages` based on `org_pkid`
UPDATE "pages" SET "space_pkid" = spaces.pkid
FROM spaces
WHERE spaces.org_pkid = pages.org_pkid;

-- 5. Remove the `org_pkid` column and foreign key constraint from `pages`
ALTER TABLE IF EXISTS "pages"
DROP CONSTRAINT IF EXISTS fk_page_organization;

ALTER TABLE IF EXISTS "pages"
DROP COLUMN IF EXISTS "org_pkid";