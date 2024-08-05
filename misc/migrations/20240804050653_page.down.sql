-- Migration down file for dropping the 'pages' table and index

DROP INDEX IF EXISTS idx_parent_page_pkid;

DROP TABLE IF EXISTS "pages";
