package model

import (
	"go-production/app/filters"
	"go-production/app/helpers"
	"go-production/global"
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

func (m Movie) List(f *filters.MovieFilter, p *helpers.Pagination, s *helpers.Ordering) ([]Movie, helpers.Metadata, error) {
	var movies []Movie
	var totalRecords int64
	var metadata helpers.Metadata

	// 2. Xây dựng câu lệnh query với filters
	query := global.DB.Model(&Movie{})
	query = f.Apply(query)

	// 3. Đếm tổng số bản ghi sau khi lọc
	if err := query.Count(&totalRecords).Error; err != nil {
		return movies, metadata, err
	}

	// 4. Truy vấn dữ liệu với phân trang
	if err := query.Limit(*p.PageSize).Offset(p.Offset).Order(s.Order).Find(&movies).Error; err != nil {
		return movies, metadata, err
	}

	p.CalculateMetadata(totalRecords)

	return movies, p.Metadata, nil

}

func (m Movie) Get(id string) (Movie, error) {
	var movie Movie
	if err := global.DB.Where("id = ?", id).First(&movie).Error; err != nil {
		return movie, err
	}
	return movie, nil
}
