package helpers

import (
	"slices"

	"github.com/gin-gonic/gin"
)

type Ordering struct {
	Order         string
	SafeOrderList []string
}

func (o *Ordering) GetParams(c *gin.Context, v *Validator) {
	sortRaw := c.Query("ordering")
	if sortRaw == "" {
		o.Order = "id DESC"
		return
	}

	var sortColumn string
	if sortRaw[0] == '-' {
		sortColumn = sortRaw[1:]
		o.Order = sortColumn + " DESC"
	} else {
		sortColumn = sortRaw
		o.Order = sortColumn + " ASC"
	}

	v.Check(slices.Contains(o.SafeOrderList, sortColumn), "ordering", "invalid sorting column")

}
