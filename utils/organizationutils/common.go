package organizationutils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const OrgPkIDParam = "orgPkID"
const OrgSlugParam = "orgSlug"

func GetOrgPkIDParam(c *gin.Context) (int64, bool) {
	orgPkID := c.Params.ByName(OrgPkIDParam)
	if orgPkID == "" {
		return int64(-1), false
	}
	docPkID, _ := strconv.Atoi(orgPkID)
	return int64(docPkID), true
}

func GetOrgSlugParam(c *gin.Context) (string, bool) {
	orgSlug := c.Params.ByName("orgSlug")
	if orgSlug == "" {
		return "", false
	}
	return orgSlug, true
}
