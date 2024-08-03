package space

type CreateSpaceDto struct {
	OwnerPkID   int64
	OrgPkID     int64
	Name        string
	Description string
}
