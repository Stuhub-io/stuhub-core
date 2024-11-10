ALTER TABLE IF EXISTS "pages"
ADD COLUMN "path" TEXT NOT NULL DEFAULT '';

-- Update paths for all outer pages
UPDATE pages
SET path = '/'
WHERE parent_page_pkid IS NULL;

-- Step 2: Recursive CTE to update paths of child pages based on their parent
WITH RECURSIVE PagePaths AS (
    -- Select root pages as the starting point
    SELECT 
        pkid, 
        parent_page_pkid, 
        path AS parent_path
    FROM 
        pages
    WHERE 
        parent_page_pkid IS NULL

    UNION ALL

    -- Recursive step: select children and append pkid to parent path
    SELECT 
        child.pkid, 
        child.parent_page_pkid, 
        CONCAT(pp.parent_path, child.pkid::text, '/') AS parent_path
    FROM 
        pages AS child
    INNER JOIN 
        PagePaths AS pp 
        ON child.parent_page_pkid = pp.pkid
)
-- Update the path field for each page based on the recursive CTE
UPDATE pages AS p
SET path = pp.parent_path
FROM PagePaths AS pp
WHERE p.pkid = pp.pkid;