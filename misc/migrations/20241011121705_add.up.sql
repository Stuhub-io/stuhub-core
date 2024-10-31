ALTER TABLE "pages"
ADD COLUMN "node_id" UUID UNIQUE;

CREATE UNIQUE INDEX page_node_id_idx ON "pages" (node_id);