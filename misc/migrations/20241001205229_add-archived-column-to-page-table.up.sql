ALTER TABLE pages
ADD  COLUMN IF NOT EXISTS archived_at TIMESTAMP WITH TIME ZONE;
