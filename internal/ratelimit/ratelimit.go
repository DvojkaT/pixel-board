package ratelimit

import (
	"sync"
	"time"
)

type User struct {
	LastPaint time.Time
}
type Limiter struct {
	Users map[uint64]*User
	mtx   sync.RWMutex
}

func NewLimiter() *Limiter {
	return &Limiter{
		Users: make(map[uint64]*User),
	}
}

func (l *Limiter) TryPaint(userId uint64) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	user, ok := l.Users[userId]
	if !ok {
		user = &User{}
		l.Users[userId] = user
	}
	if time.Since(user.LastPaint) > 500*time.Millisecond {
		user.LastPaint = time.Now()
		return true
	}
	return false
}
