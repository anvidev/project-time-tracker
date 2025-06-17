package categories

import (
	"context"
	"errors"
	"strings"
)

var (
	ErrAlreadyFollowed      = errors.New("already following category")
	ErrNotFollowingCategory = errors.New("category is not followed")
	ErrCategoryNotFollowed  = errors.New("category was not followed")
)

func (s *Store) Leafs(ctx context.Context, userId int64) ([]Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	categories := []Category{}

	stmt := `
		with parent(id, title, parent_id, is_retired, path_parent_retired, root_parent_title) as (
		  select id, title, parent_id, is_retired, is_retired as path_parent_retired, title as root_parent_title
		  from categories c
          join users_categories_link ucl on ucl.category_id = c.id
		  where ucl.user_id = ?
          
		  union all
          
		  select c.id, c.title, c.parent_id, c.is_retired, p.path_parent_retired or c.is_retired, p.root_parent_title
		  from categories c
		  join parent p on c.parent_id = p.id
		),
		non_leafs as (
		  select distinct parent_id as id
		  from categories
		  where parent_id is not null
		)
		select distinct c.id, c.title, p.root_parent_title
		from categories c
		join parent p on c.id = p.id
		where c.id not in (select id from non_leafs)
		  and p.path_parent_retired = 0;
	`

	rows, err := s.db.QueryContext(ctx, stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c Category
		rows.Scan(&c.Id, &c.Title, &c.RootTitle)
		categories = append(categories, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *Store) Follow(ctx context.Context, id, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		insert into users_categories_link (category_id, user_id)
		values (?, ?)
	`

	result, err := s.db.ExecContext(ctx, stmt, id, userId)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "UNIQUE constraint failed"):
			return ErrAlreadyFollowed
		default:
			return err
		}
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrCategoryNotFollowed
	}

	return nil
}

func (s *Store) Unfollow(ctx context.Context, id, userId int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `delete from users_categories_link where category_id = ? and user_id = ?`

	result, err := s.db.ExecContext(ctx, stmt, id, userId)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrNotFollowingCategory
	}

	return nil
}
