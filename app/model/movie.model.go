package model

import (
	"context"
	"errors"
	"go-production/app/filters"
	"go-production/app/helpers"
	"go-production/global"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
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

type InputUpdateMovie struct {
	Title   *string  `json:"title"`
	Year    *int32   `json:"year"`
	Runtime *int32   `json:"runtime"`
	Genres  []string `json:"genres"`
}

// TableName chỉ định tên bảng cho GORM (mặc định là "movies")
func (Movie) TableName() string {
	return "movies"
}

func (m Movie) List(f *filters.MovieFilter, p *helpers.Pagination, s *helpers.Ordering) ([]Movie, helpers.Metadata, error) {
	var movies []Movie
	var metadata helpers.Metadata

	// 2. Xây dựng câu lệnh query với filters
	query := global.DB.Model(&Movie{})
	query = f.Apply(query)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 3. Áp dụng Keyset Pagination logic
	if p.LastID > 0 {
		query = query.Where("id > ?", p.LastID)
	}

	// 4. Truy vấn dữ liệu với chiến thuật Limit + 1
	order := s.Order
	if order == "" {
		order = "id ASC"
	}

	// Lấy dư 1 bản ghi để kiểm tra xem còn trang sau không
	limit := *p.PageSize + 1
	if err := query.WithContext(ctx).Limit(limit).Order(order).Find(&movies).Error; err != nil {
		return nil, metadata, err
	}

	// 5. Kiểm tra HasNext và xử lý kết quả
	hasNext := false
	var nextCursor int64

	if len(movies) > *p.PageSize {
		hasNext = true
		movies = movies[:*p.PageSize] // Cắt bỏ bản ghi thứ n+1
	}

	if len(movies) > 0 {
		nextCursor = movies[len(movies)-1].ID
	}

	p.CalculateMetadata(hasNext, nextCursor)

	return movies, p.Metadata, nil
}
}

func (m Movie) Get(id string) (Movie, error) {
	var movie Movie
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := global.DB.WithContext(ctx).Where("id = ?", id).First(&movie).Error; err != nil {
		return movie, err
	}
	return movie, nil
}

func (m Movie) Create(movie *Movie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := global.DB.WithContext(ctx).Create(movie).Error; err != nil {
		return err
	}
	return nil
}

func (m *Movie) Update(input *InputUpdateMovie) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if input.Title != nil {
		m.Title = *input.Title
	}
	if input.Year != nil {
		m.Year = *input.Year
	}
	if input.Runtime != nil {
		m.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		m.Genres = input.Genres
	}

	result := global.DB.WithContext(ctx).
		Model(m).
		Where("id = ? AND version = ?", m.ID, m.Version).
		Updates(map[string]interface{}{
			"title":   m.Title,
			"year":    m.Year,
			"runtime": m.Runtime,
			"genres":  m.Genres,
			"version": gorm.Expr("version + 1"),
		})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("edit conflict")
	}

	m.Version++
	return nil
}
