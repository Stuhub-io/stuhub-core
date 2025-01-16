ALTER TABLE page_access_logs
DROP CONSTRAINT IF EXISTS unique_page_access_log,
ADD CONSTRAINT unique_page_access_log UNIQUE (page_pkid, user_pkid);