package sessions

import "time"

type Session struct {
	Token     string `json:"token" apiduck:"desc=Session is valid for 7 days and is extended on each request"`
	UserId    int64  `json:"userId"`
	ExpiresAt string `json:"expiresAt"` // yyyy-MM-dd (time.DateOnly)
	CreatedAt string `json:"createdAt"` // yyyy-MM-dd (time.DateOnly)
	UpdatedAt string `json:"updatedAt"` // yyyy-MM-dd (time.DateOnly)
}

func (s Session) IsExpired() bool {
	expires, err := time.Parse(time.DateTime, s.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expires)
}
