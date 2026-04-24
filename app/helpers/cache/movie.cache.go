package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"go-production/global"
	"sort"
	"strings"
	"time"
)

const movieTagKey = "movies:tag:list"

func BuildMovieKey(title string, genres string, pageSize int, lastID int) string {
	//Sort genres để đảm bảo key cố định
	genresSplitted := strings.Split(genres, ",")
	sort.Strings(genresSplitted)
	sortedGenres := strings.Join(genresSplitted, ",")

	key := fmt.Sprintf(
		"movies:data:ps:%d:lid:%d:title:%s:genres:%s",
		pageSize,
		lastID,
		strings.ToLower(strings.TrimSpace(title)),
		strings.ToLower(strings.TrimSpace(sortedGenres)),
	)

	hash := md5.Sum([]byte(key))
	return fmt.Sprintf("movies:data:%x", hash)
}

func GetMovieCache(ctx context.Context, key string) (string, bool) {
	val, err := global.Redis.Get(ctx, key).Result()
	if err != nil {
		return "", false
	}
	return val, true
}

func SetMovieCache(ctx context.Context, key string, data string, ttl time.Duration) error {
	pipe := global.Redis.Pipeline()
	// 1. Lưu data (JSON string) với TTL
	pipe.Set(ctx, key, data, ttl)
	// 2. Đăng ký key vào Tag Set để invalidation biết cần xóa key nào
	//    Tag Set không cần TTL vì ta sẽ DEL nó khi invalidate.
	//    Nhưng đặt TTL dài hơn data key để tránh orphan entries.
	pipe.SAdd(ctx, movieTagKey, key)
	pipe.Expire(ctx, movieTagKey, ttl+1*time.Minute)
	_, err := pipe.Exec(ctx)
	return err
}

func InvalidateMovieCache(ctx context.Context) error {
	// Lấy tất cả keys đã đăng ký trong tag set
	keys, err := global.Redis.SMembers(ctx, movieTagKey).Result()
	if err != nil || len(keys) == 0 {
		// Nếu tag set trống hoặc lỗi, vẫn xóa tag set cho sạch
		global.Redis.Del(ctx, movieTagKey)
		return nil
	}
	pipe := global.Redis.Pipeline()
	// UNLINK = async DEL, không block Redis event loop
	// Khác với DEL: Redis đánh dấu keys để xóa ở background goroutine
	pipe.Unlink(ctx, keys...)
	// Xóa tag set để reset danh sách
	pipe.Del(ctx, movieTagKey)
	_, err = pipe.Exec(ctx)
	return err
}
