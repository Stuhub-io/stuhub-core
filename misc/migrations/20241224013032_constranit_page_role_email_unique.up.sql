CREATE UNIQUE INDEX IF NOT EXISTS page_roles_page_email_unique_idx 
	ON "page_roles" ("page_pkid", "email");