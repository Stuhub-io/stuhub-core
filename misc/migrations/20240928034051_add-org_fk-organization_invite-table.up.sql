ALTER TABLE "organization_invites"
ADD CONSTRAINT fk_organization
FOREIGN KEY (organization_pkid)
REFERENCES "organizations" (pkid)
ON DELETE CASCADE;