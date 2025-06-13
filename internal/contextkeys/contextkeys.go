package contextkeys

type contextkey string

var (
	SessionToken contextkey = "session_id"
	UserId                  = "user_id"
)
