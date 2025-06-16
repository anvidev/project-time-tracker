-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
  token text primary key,
  user_id integer not null references users (id),
  expires_at text not null,
  created_at text not null,
  updated_at text not null
);

create index idx_sessions_user_id on sessions (user_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop index if exists idx_sessions_user_id;

drop table if exists sessions;

-- +goose StatementEnd
