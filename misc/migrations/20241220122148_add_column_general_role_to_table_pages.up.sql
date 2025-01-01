ALTER TABLE "pages"
ADD COLUMN IF NOT EXISTS "general_role" VARCHAR(20) DEFAULT 'viewer',
DROP CONSTRAINT IF EXISTS "check_page_general_role",
ADD CONSTRAINT  "check_page_general_role" CHECK ("general_role" IN ('viewer', 'editor'));