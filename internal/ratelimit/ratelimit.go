package ratelimit

import (
	"sync"
	"time"
)

type User struct {
	LastPaint time.Time
}
type Limiter struct {
	Users map[string]*User
	mtx   sync.RWMutex
}

func NewLimiter() *Limiter {
	return &Limiter{
		Users: make(map[string]*User),
	}
}

func (l *Limiter) TryPaint(userId string) bool {
	l.mtx.Lock()
	defer l.mtx.Unlock()
	user, ok := l.Users[userId]
	if !ok {
		user = &User{}
		l.Users[userId] = user
	}
	if time.Since(user.LastPaint) > 100*time.Millisecond {
		user.LastPaint = time.Now()
		return true
	}
	return false
}
