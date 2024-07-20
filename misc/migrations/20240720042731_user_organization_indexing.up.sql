-- Create indexes for the 'organization_member' table
CREATE INDEX idx_organization_member_organization_pkid ON "organization_member" (organization_pkid);
CREATE INDEX idx_organization_member_user_pkid ON "organization_member" (user_pkid);