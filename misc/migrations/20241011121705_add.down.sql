ALTER TABLE "pages"
DROP COLUMN "node_id";

DROP INDEX IF EXISTS page_node_id_idx;
