package web

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"sumeru/core/server/config"
)

type rateBucket struct {
	count  int
	window time.Time
}

var (
	rateMu       sync.Mutex
	rateByIP     = map[string]*rateBucket{}
	rateLimitOn  bool
	rateLimitRPM int
)

// InitRateLimit reads rate_limit_rpm from loaded config.
func InitRateLimit() {
	rateLimitRPM = config.AppConfig.RateLimitRPM
	rateLimitOn = rateLimitRPM > 0
}

func rateLimitedPath(path string) bool {
	switch path {
	case apiRPCRoute, loginRoute, totpLoginRoute, oauthStartRoute, oauthCallbackRoute:
		return true
	default:
		return false
	}
}

func allowRate(clientIP string) bool {
	if !rateLimitOn {
		return true
	}
	clientIP = strings.TrimSpace(clientIP)
	if clientIP == "" {
		clientIP = "unknown"
	}
	now := time.Now()
	rateMu.Lock()
	defer rateMu.Unlock()
	b, ok := rateByIP[clientIP]
	if !ok || now.Sub(b.window) >= time.Minute {
		rateByIP[clientIP] = &rateBucket{count: 1, window: now}
		return true
	}
	if b.count >= rateLimitRPM {
		return false
	}
	b.count++
	return true
}

func enforceRateLimit(w http.ResponseWriter, r *http.Request) bool {
	if !rateLimitedPath(r.URL.Path) {
		return true
	}
	if allowRate(clientIP(r)) {
		return true
	}
	http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
	return false
}

// pruneAttemptsAfter keeps timestamps strictly after cutoff (in-place slice reuse).
func pruneAttemptsAfter(attempts []time.Time, cutoff time.Time) []time.Time {
	out := attempts[:0]
	for _, attemptTime := range attempts {
		if attemptTime.After(cutoff) {
			out = append(out, attemptTime)
		}
	}
	return out
}

// pruneAttemptsWithin keeps timestamps within window ending at now.
func pruneAttemptsWithin(attempts []time.Time, now time.Time, window time.Duration) []time.Time {
	cutoff := now.Add(-window)
	return pruneAttemptsAfter(attempts, cutoff)
}
