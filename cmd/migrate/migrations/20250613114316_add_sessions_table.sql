-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
  token text primary key,
  user_id integer not null references users (id),
  expires_at text not null,
  created_at text not null,
  updated_at text not null
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists sessions;

-- +goose StatementEnd
