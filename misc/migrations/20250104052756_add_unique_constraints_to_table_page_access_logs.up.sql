ALTER TABLE page_access_logs
ADD CONSTRAINT unique_page_access_log UNIQUE (page_pkid, user_pkid);