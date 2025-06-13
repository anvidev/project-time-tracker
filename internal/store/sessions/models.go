package sessions

type Session struct {
	Token     int64  `json:"token"`
	UserId    int64  `json:"userId"`
	ExpiresAt string `json:"expiresAt"`
}
