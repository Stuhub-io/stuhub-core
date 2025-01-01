ALTER TABLE pages
ADD COLUMN IF NOT EXISTs author_pkid BIGINT;

-- add foreign key constraint
ALTER TABLE pages DROP CONSTRAINT IF EXISTS fk_page_author;

ALTER TABLE pages
ADD CONSTRAINT fk_page_author
FOREIGN KEY (author_pkid) REFERENCES users(pkid)
ON DELETE SET NULL; 