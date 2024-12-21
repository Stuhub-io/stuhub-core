ALTER TABLE "pages"
ADD COLUMN "general_role" VARCHAR(20) DEFAULT 'viewer',
ADD CONSTRAINT "check_page_general_role" CHECK ("general_role" IN ('viewer', 'editor'));