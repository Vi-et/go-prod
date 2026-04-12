package middleware

import (
	"expvar"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	totalRequests         = expvar.NewInt("total_requests_count")
	totalResponses        = expvar.NewInt("total_responses_count")
	totalProcessingTimeMs = expvar.NewInt("total_processing_time_ms")
	responsesMap          = expvar.NewMap("total_responses_by_status")
)

func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		totalRequests.Add(1)

		c.Next()

		duration := time.Since(start).Milliseconds()
		totalProcessingTimeMs.Add(duration)
		totalResponses.Add(1)

		status := strconv.Itoa(c.Writer.Status())
		responsesMap.Add(status, 1)
	}
}
