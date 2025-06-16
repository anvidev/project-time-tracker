-- +goose Up
-- +goose StatementBegin
create table if not exists users_categories_link (
  user_id integer not null references users (id),
  category_id integer not null references categories (id),
  primary key (user_id, category_id)
);

create index idx_users_categories_category_id on users_categories_link (category_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop index if exists idx_users_categories_category_id;

drop table if exists users_categories_link;

-- +goose StatementEnd
