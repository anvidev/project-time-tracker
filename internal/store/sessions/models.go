package sessions

import "time"

type Session struct {
	Token     string `json:"token"`
	UserId    int64  `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

func (s Session) IsExpired() bool {
	expires, err := time.Parse(time.DateTime, s.ExpiresAt)
	if err != nil {
		return true
	}
	return time.Now().After(expires)
}
