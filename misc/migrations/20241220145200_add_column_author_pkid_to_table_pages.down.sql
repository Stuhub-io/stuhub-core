ALTER TABLE pages
DROP CONSTRAINT fk_page_author;

ALTER TABLE pages
DROP COLUMN author_pkid;
