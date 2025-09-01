package categories

import (
	"context"
	"database/sql"
	"errors"
	"strings"
)

var (
	ErrAlreadyFollowed      = errors.New("already following category")
	ErrNotFollowingCategory = errors.New("category is not followed")
	ErrCategoryNotFollowed  = errors.New("category was not followed")
	ErrCategoryNotToggled   = errors.New("category was not toggled")
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

func (s *Store) Tree(ctx context.Context, userId int64) ([]*CategoryTree, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		select 
		  c.id,
		  c.parent_id,
		  c.title,
		  c.is_retired,
		  (select exists(select 1 from users_categories_link where user_id = ? and category_id = c.id)) as is_followed
		from categories c
		order by c.parent_id nulls first, c.id
	`

	rows, err := s.db.QueryContext(ctx, stmt, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allCategories := make(map[int64]*CategoryTree)
	var tree []*CategoryTree

	for rows.Next() {
		var category CategoryTree

		if err := rows.Scan(
			&category.Id,
			&category.ParentId,
			&category.Title,
			&category.IsRetired,
			&category.IsFollowed,
		); err != nil {
			return nil, err
		}

		category.Children = make([]*CategoryTree, 0)
		allCategories[category.Id] = &category
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	for _, category := range allCategories {
		if category.ParentId == nil {
			tree = append(tree, category)
		} else {
			if parent, exists := allCategories[*category.ParentId]; exists {
				parent.Children = append(parent.Children, category)
			}
		}
	}

	return tree, nil
}

func (s *Store) Create(ctx context.Context, input CreateCategoryInput) (*Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		insert into categories (title, parent_id)
		values (?, ?)
		returning id, coalesce((select title from categories where id = ?), '') as root_title
	`

	category := Category{Title: input.Title}
	var rootTitle sql.NullString

	err := s.db.
		QueryRowContext(ctx, stmt, input.Title, input.ParentId, input.ParentId).
		Scan(&category.Id, &rootTitle)
	if err != nil {
		return nil, err
	}

	category.RootTitle = rootTitle.String

	return &category, nil
}

func (s *Store) Update(ctx context.Context, id int64, title string) (*Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		update categories c
		set title = ?
		where id = ? 
		returning id, coalesce((select title from categories where id = c.parent_id), '') as root_title
	`

	category := Category{Title: title}
	var rootTitle sql.NullString

	err := s.db.QueryRowContext(ctx, stmt, title, id).Scan(&category.Id, &rootTitle)
	if err != nil {
		return nil, err
	}

	category.RootTitle = rootTitle.String

	return &category, nil

}

func (s *Store) ToggleRetire(ctx context.Context, id int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `update categories set is_retired = not is_retired where id = ?`

	result, err := s.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected != 1 {
		return ErrCategoryNotToggled
	}

	return nil
}

func (s *Store) Get(ctx context.Context, id int64) (*Category, error) {
	ctx, cancel := context.WithTimeout(ctx, s.queryTimeout)
	defer cancel()

	stmt := `
		select
			c.id,
			c.title,
			coalesce((select title from categories where id = c.parent_id), '') as root_title
		from categories c
		where c.id = ?
	`

	var c Category

	if err := s.db.QueryRowContext(ctx, stmt, id, id).Scan(&c.Id, &c.Title, &c.RootTitle); err != nil {
		return nil, err
	}

	return &c, nil
}
