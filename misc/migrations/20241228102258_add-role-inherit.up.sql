-- Add role value "inherit" 
ALTER TABLE "pages"
DROP CONSTRAINT IF EXISTS "check_page_general_role",
ADD CONSTRAINT "check_page_general_role" CHECK ("general_role" IN ('viewer', 'editor', 'inherit'));

ALTER TABLE "page_roles"
DROP CONSTRAINT IF EXISTS "page_roles_role_check",
ADD CONSTRAINT "page_roles_role_check" CHECK ("role" IN ('viewer', 'editor', 'inherit'));