-- +goose Up
-- +goose StatementBegin
create table if not exists users (
  id integer primary key,
  name text not null,
  email text unique not null,
  hash blob not null,
  role text not null,
  is_active integer not null default 1,
  created_at text not null
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists users;

-- +goose StatementEnd
