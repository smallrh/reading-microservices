package observability

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"time"
)

var (
	UserRegistrations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "user_registrations_total",
		Help: "Total number of user registrations",
	})

	UserLogins = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "user_logins_total",
		Help: "Total number of user logins",
	}, []string{"status"})

	DatabaseQueryDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "database_query_duration_seconds",
		Help:    "Database query duration in seconds",
		Buckets: []float64{0.001, 0.005, 0.01, 0.05, 0.1, 0.5},
	}, []string{"operation"})
)

func init() {
	prometheus.MustRegister(UserRegistrations, UserLogins, DatabaseQueryDuration)
}

func LoggingMiddleware(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		logger.WithFields(logrus.Fields{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"status":     c.Writer.Status(),
			"duration":   time.Since(start),
			"client_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		}).Info("HTTP request")
	}
}
