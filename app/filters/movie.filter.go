package filters

import (
	"go-production/app/helpers"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type MovieFilter struct {
	Title      string
	Genres     string
	Year       *int
	YearGTE    *int
	YearLTE    *int
	Runtime    *int
	RuntimeGTE *int
	RuntimeLTE *int
}

func (f *MovieFilter) GetParams(c *gin.Context, v *helpers.Validator) {
	f.Title = c.Query("title")
	f.Genres = c.Query("genres")
	f.Year = helpers.GetIntParam(c, v, "year")
	f.YearGTE = helpers.GetIntParam(c, v, "yearGTE")
	f.YearLTE = helpers.GetIntParam(c, v, "yearLTE")
	f.Runtime = helpers.GetIntParam(c, v, "runtime")
	f.RuntimeGTE = helpers.GetIntParam(c, v, "runtimeGTE")
	f.RuntimeLTE = helpers.GetIntParam(c, v, "runtimeLTE")
}

func (f *MovieFilter) Apply(db *gorm.DB) *gorm.DB {
	if f.Title != "" {
		db = db.Where("to_tsvector('simple', title) @@ plainto_tsquery('simple', ?)", f.Title)
	}
	if f.Genres != "" {
		db = db.Where("genres @> ?", pq.Array(f.Genres))
	}
	if f.Year != nil {
		db = db.Where("year = ?", f.Year)
	}
	if f.YearGTE != nil {
		db = db.Where("year >= ?", f.YearGTE)
	}
	if f.YearLTE != nil {
		db = db.Where("year <= ?", f.YearLTE)
	}
	if f.Runtime != nil {
		db = db.Where("runtime = ?", f.Runtime)
	}
	if f.RuntimeGTE != nil {
		db = db.Where("runtime >= ?", f.RuntimeGTE)
	}
	if f.RuntimeLTE != nil {
		db = db.Where("runtime <= ?", f.RuntimeLTE)
	}
	return db
}
