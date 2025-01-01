ALTER TABLE "pages" DROP COLUMN IF EXISTS "is_general_access",
DROP CONSTRAINT IF EXISTS "check_page_general_role",
ADD CONSTRAINT "check_page_general_role" CHECK ("general_role" IN ('viewer', 'editor', 'inherit', 'restricted'));