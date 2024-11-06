ALTER TABLE "pages"
ADD COLUMN "org_pkid" BIGINT;

ALTER TABLE "pages"
    ADD CONSTRAINT fk_page_organization
    FOREIGN KEY (org_pkid)
    REFERENCES "organizations" (pkid)
    ON DELETE CASCADE;


UPDATE "pages" SET "org_pkid" = spaces.org_pkid
FROM spaces
WHERE spaces.pkid = pages.space_pkid;

ALTER TABLE "pages"
DROP COLUMN "space_pkid";

DROP TABLE IF EXISTS "space_member";
DROP TABLE IF EXISTS "spaces";