package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetIntParam(c *gin.Context, v *Validator, key string) *int {
	val := c.Query(key)

	if val == "" {
		return nil // Trả về nil nếu không truyền
	}

	num, err := strconv.Atoi(val)
	if err != nil {
		v.AddError(key, "invalid_integer")
		return nil
	}

	return &num // Trả về địa chỉ của biến num
}
