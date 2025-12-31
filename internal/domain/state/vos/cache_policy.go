package vos

import "time"

const defaultTTL = 30 * 24 * time.Hour // 30 días

type CachePolicy struct {
	ttl time.Duration
}

func NewCachePolicy(ttl time.Duration) CachePolicy {
	if ttl <= 0  {
		ttl = defaultTTL
	}
	return CachePolicy{ttl: ttl}
}

func (p CachePolicy) TTL() time.Duration {
	return p.ttl
}
