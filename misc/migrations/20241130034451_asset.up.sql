-- Drop the existing CHECK constraint
ALTER TABLE IF EXISTS "pages" DROP CONSTRAINT IF EXISTS "pages_view_type_check";

-- Add the new CHECK constraint with the additional option
ALTER TABLE IF EXISTS "pages" 
ADD CONSTRAINT "pages_view_type_check" 
CHECK (view_type IN ('document', 'folder', 'asset'));

CREATE TABLE IF NOT EXISTS "assets" (
    "pkid" bigserial PRIMARY KEY,
    "page_pkid" BIGINT NOT NULL UNIQUE,
    "url" TEXT NOT NULL DEFAULT '',
    "size" BIGINT,
    "extension" CHAR(100),
    "updated_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "thumbnails" JSONB NOT NULL DEFAULT '{}'::JSONB,

    CONSTRAINT fk_asset_page
        FOREIGN KEY (page_pkid)
        REFERENCES "pages" (pkid)
        ON DELETE CASCADE
)