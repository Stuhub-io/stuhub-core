package pageutils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	PagePkIDParam = "pagePkID"
	PageIDParam   = "pageID"
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

func AppendPath(path string, id string) string {
	if path == "" {
		return id
	}
	return path + "/" + id
}
