BEGIN;

-- Drop the existing CHECK constraint
ALTER TABLE "pages" 
DROP CONSTRAINT IF EXISTS "pages_view_type_check";

-- Add the new CHECK constraint with updated values
ALTER TABLE "pages" 
ADD CONSTRAINT "pages_view_type_check" 
CHECK (view_type IN ('document', 'folder'));

COMMIT;