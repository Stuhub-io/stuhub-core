package pageutils

import (
	"strconv"
	"strings"

	sliceutils "github.com/Stuhub-io/utils/slice"
	"github.com/gin-gonic/gin"
)

const (
	PagePkIDParam      = "pagePkID"
	PageIDParam        = "pageID"
	PublicTokenIDParam = "publicTokenID"
)

func GetPageIDParam(c *gin.Context) (string, bool) {
	pageID := c.Params.ByName(PageIDParam)
	if pageID == "" {
		return "", false
	}
	return pageID, true
}

func GetPagePkIDParam(c *gin.Context) (int64, bool) {
	pagePkID := c.Params.ByName(PagePkIDParam)
	if pagePkID == "" {
		return int64(-1), false
	}
	docPkID, _ := strconv.Atoi(pagePkID)
	return int64(docPkID), true
}

func GetPublicTokenIDParam(c *gin.Context) (string, bool) {
	publicTokenID := c.Params.ByName(PublicTokenIDParam)
	if publicTokenID == "" {
		return "", false
	}
	return publicTokenID, true
}

func AppendPath(path string, pkID string) string {
	if path == "" {
		return pkID
	}
	return path + "/" + pkID
}

func PagePathToPkIDs(path string) []int64 {
	parentPkIDs := sliceutils.Map(strings.Split(path, "/"), func(pkid string) int64 {
		parsedPkID, err := strconv.ParseInt(pkid, 10, 64)
		if err != nil {
			return -1
		}
		return parsedPkID
	})
	return sliceutils.Filter(parentPkIDs, func(pkid int64) bool {
		return pkid != -1
	})
}

func BuildPagePath(pkIDs []int64) string {
	if len(pkIDs) == 0 {
		return ""
	}
	path := strconv.FormatInt(pkIDs[0], 10)
	for _, pkID := range pkIDs[1:] {
		path += "/" + strconv.FormatInt(pkID, 10)
	}
	return path
}
