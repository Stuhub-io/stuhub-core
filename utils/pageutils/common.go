package pageutils

import (
	"strconv"
	"strings"

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

func AppendPath(path string, id string) string {
	if path == "" {
		return id
	}
	return path + "/" + id
}

func PagePathToPkIDs(path string) []int64 {
	if len(path) == 0 {
		return []int64{}
	}
	pkIdStrs := strings.Split(path, "/")
	pkIDs := make([]int64, 0, len(pkIdStrs))
	for _, pkIdStr := range pkIdStrs {
		pkID, _ := strconv.Atoi(pkIdStr)
		pkIDs = append(pkIDs, int64(pkID))
	}
	return pkIDs
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
