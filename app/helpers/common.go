package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIntParam(c *gin.Context, v *Validator, key string) int {
	val := c.Query(key)
	num, err := strconv.Atoi(val)
	if err != nil {
		v.AddError(key, "invalid_integer")
		return 0
	}
	return num
}
