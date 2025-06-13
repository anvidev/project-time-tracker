package sessions

import (
	"context"
	"errors"
	"time"

	"github.com/anvidev/project-time-tracker/internal/id"
)

var (
	ErrConflictNoUser    = errors.New("user id not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrSessionNotCreated = errors.New("session not created")
	ErrSessionNotDeleted = errors.New("session not deleted")
	ErrSessionExpired    = errors.New("session expired")
)

func (s *Store) Create(ctx context.Context, userId int64) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	session := &Session{
		Token:     id.String(32, id.Numbers, id.LettersUpper),
		UserId:    userId,
		ExpiresAt: time.Now().Add(s.sessionExpiresIn).Format(time.DateTime),
		CreatedAt: time.Now().Format(time.DateTime),
		UpdatedAt: time.Now().Format(time.DateTime),
	}

	stmt := `
		insert into sessions (token, user_id, expires_at, created_at, updated_at)
		values (?, ?, ?, ?, ?)
	`

	result, err := s.db.ExecContext(
		ctx,
		stmt,
		session.Token,
		session.UserId,
		session.ExpiresAt,
		session.CreatedAt,
		session.UpdatedAt,
	)
	if err != nil {
		switch {
		case err.Error() == "unique something":
			return nil, ErrConflictNoUser
		default:
			return nil, err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if affected != 1 {
		return nil, ErrSessionNotCreated
	}

	return session, nil
}

func (s *Store) Validate(ctx context.Context, token string) (*Session, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	var session Session

	stmt := `
		select token, user_id, expires_at, created_at, updated_at
		from sessions
		where token = ?
	`

	if err := s.db.QueryRowContext(
		ctx,
		stmt,
		token,
	).Scan(
		&session.Token,
		&session.UserId,
		&session.ExpiresAt,
		&session.CreatedAt,
		&session.UpdatedAt,
	); err != nil {
		return nil, err
	}

	if session.IsExpired() {
		stmt = `delete from sessions where token = ?`

		result, err := s.db.ExecContext(ctx, stmt, token)
		if err != nil {
			return nil, err
		}

		affected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}

		if affected != 1 {
			return nil, ErrSessionNotDeleted
		}

		return nil, ErrSessionExpired
	} else {
		stmt = `
			update sessions
			set expires_at = ?, updated_at = ?
			where token = ? 
		`

		session.ExpiresAt = time.Now().Add(s.sessionExpiresIn).Format(time.DateTime)
		session.UpdatedAt = time.Now().Format(time.DateTime)

		_, err := s.db.ExecContext(
			ctx,
			stmt,
			session.ExpiresAt,
			session.UpdatedAt,
			token,
		)
		if err != nil {
			return nil, err
		}

		return &session, nil
	}
}

func (s *Store) InvalidateAll(ctx context.Context, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `delete from sessions where user_id = ?`

	result, err := s.db.ExecContext(ctx, stmt, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrSessionNotDeleted
	}

	return nil
}
