package request

import "github.com/gin-gonic/gin"

func Validate(c *gin.Context, data any) error {
	err := c.Bind(data)
	return err
}
