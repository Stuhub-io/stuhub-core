ALTER TABLE "pages"
ADD COLUMN IF NOT EXISTS "is_general_access" BOOLEAN NOT NULL DEFAULT FALSE;