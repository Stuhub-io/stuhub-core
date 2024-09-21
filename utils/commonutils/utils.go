package commonutils

import "github.com/gin-gonic/gin"

func ParamsToMap(p gin.Params) map[string]string {
	params := make(map[string]string, len(p))
	for _, param := range p {
		params[param.Key] = param.Value
	}
	return params
}
