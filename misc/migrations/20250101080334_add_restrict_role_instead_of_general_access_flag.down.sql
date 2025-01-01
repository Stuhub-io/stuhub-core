ALTER TABLE "pages" ADD COLUMN "is_general_access" BOOLEAN DEFAULT FALSE;

DROP CONSTRAINT IF EXISTS "check_page_general_role",
ADD CONSTRAINT "check_page_general_role" CHECK ("general_role" IN ('viewer', 'editor', 'inherit'));
