package sessions

import (
	"context"
	"errors"
	"time"
)

var (
	ErrConflictNoUser = errors.New("user id not found")
)

func (s *Store) Create(ctx context.Context, userId int64) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	session := &Session{
		UserId:    userId,
		ExpiresAt: time.Now().Add(s.sessionExpiresIn).Format(time.DateTime),
	}

	stmt := `
		insert into sessions (user_id, expires_at)
		values (?, ?)
		returning rowid
	`

	if err := s.db.QueryRowContext(
		ctx,
		stmt,
		session.UserId,
		session.ExpiresAt,
	).Scan(
		&session.Token,
	); err != nil {
		switch {
		case err.Error() == "unique something":
			return nil, ErrConflictNoUser
		default:
			return nil, err
		}
	}

	return session, nil
}

func (s *Store) Validate(ctx context.Context, token int64) (*Session, error) {
	return nil, nil
}
func (s *Store) Invalidate(ctx context.Context, token int64) error {
	return nil
}
