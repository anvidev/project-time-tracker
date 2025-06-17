package categories

import "context"

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
