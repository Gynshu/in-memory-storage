package domain

type RateLimiter interface {
	Limit(ip string)
	Check(ip string) bool
}
