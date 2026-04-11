package model

import (
	"time"

	"github.com/lib/pq"
)

type Movie struct {
	ID        int64          `json:"id"`
	CreatedAt time.Time      `json:"-"`
	Title     string         `json:"title"`
	Year      int32          `json:"year,omitempty"`
	Runtime   int32          `json:"runtime,omitempty"`
	Genres    pq.StringArray `json:"genres,omitempty" gorm:"type:text[]"`
	Version   int32          `json:"version"`
}

// TableName chỉ định tên bảng cho GORM (mặc định là "movies")
func (Movie) TableName() string {
	return "movies"
}
