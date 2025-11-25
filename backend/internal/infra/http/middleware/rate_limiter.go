package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
)

// RateLimiterConfig configura o rate limiter
type RateLimiterConfig struct {
	RequestsPerWindow int           // Número máximo de requests
	Window            time.Duration // Janela de tempo
	KeyFunc           func(c echo.Context) string
}

// RateLimiter implementa rate limiting em memória
type RateLimiter struct {
	config  RateLimiterConfig
	clients map[string]*clientBucket
	mu      sync.RWMutex
}

type clientBucket struct {
	requests  int
	resetTime time.Time
}

// NewRateLimiter cria um novo rate limiter
func NewRateLimiter(config RateLimiterConfig) *RateLimiter {
	if config.KeyFunc == nil {
		config.KeyFunc = func(c echo.Context) string {
			return c.RealIP()
		}
	}

	rl := &RateLimiter{
		config:  config,
		clients: make(map[string]*clientBucket),
	}

	// Limpar buckets expirados a cada minuto
	go rl.cleanup()

	return rl
}

// Middleware retorna o middleware de rate limiting
func (rl *RateLimiter) Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			key := rl.config.KeyFunc(c)

			rl.mu.Lock()
			bucket, exists := rl.clients[key]
			now := time.Now()

			if !exists || now.After(bucket.resetTime) {
				// Criar novo bucket ou resetar
				bucket = &clientBucket{
					requests:  1,
					resetTime: now.Add(rl.config.Window),
				}
				rl.clients[key] = bucket
				rl.mu.Unlock()
				return next(c)
			}

			if bucket.requests >= rl.config.RequestsPerWindow {
				rl.mu.Unlock()
				retryAfter := int(time.Until(bucket.resetTime).Seconds())
				c.Response().Header().Set("Retry-After", string(rune(retryAfter)))
				c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.RequestsPerWindow)))
				c.Response().Header().Set("X-RateLimit-Remaining", "0")
				return echo.NewHTTPError(http.StatusTooManyRequests, "Rate limit exceeded. Try again later.")
			}

			bucket.requests++
			remaining := rl.config.RequestsPerWindow - bucket.requests
			rl.mu.Unlock()

			c.Response().Header().Set("X-RateLimit-Limit", string(rune(rl.config.RequestsPerWindow)))
			c.Response().Header().Set("X-RateLimit-Remaining", string(rune(remaining)))

			return next(c)
		}
	}
}

// cleanup remove buckets expirados periodicamente
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.clients {
			if now.After(bucket.resetTime) {
				delete(rl.clients, key)
			}
		}
		rl.mu.Unlock()
	}
}

// RateLimitByUserID cria um rate limiter baseado em user_id
func RateLimitByUserID(requestsPerWindow int, window time.Duration) echo.MiddlewareFunc {
	limiter := NewRateLimiter(RateLimiterConfig{
		RequestsPerWindow: requestsPerWindow,
		Window:            window,
		KeyFunc: func(c echo.Context) string {
			userID := GetUserIDFromContext(c)
			return "user:" + userID
		},
	})
	return limiter.Middleware()
}

// RateLimitExportData rate limiter específico para export (1x por dia)
func RateLimitExportData() echo.MiddlewareFunc {
	return RateLimitByUserID(1, 24*time.Hour)
}
