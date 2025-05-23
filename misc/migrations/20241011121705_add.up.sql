ALTER TABLE IF EXISTS "pages"
ADD COLUMN IF NOT EXISTS "node_id" UUID UNIQUE;

CREATE UNIQUE INDEX IF NOT EXISTS page_node_id_idx ON "pages" (node_id);