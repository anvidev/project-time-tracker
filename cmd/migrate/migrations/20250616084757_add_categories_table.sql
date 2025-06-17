-- +goose Up
-- +goose StatementBegin
create table if not exists categories (
  id integer primary key,
  parent_id integer references categories (id) default null,
  title text not null,
  is_retired integer not null default 0
);

create index idx_categories_parent_id on categories (parent_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop index if exists idx_categories_parent_id;

drop table if exists categories;

-- +goose StatementEnd
