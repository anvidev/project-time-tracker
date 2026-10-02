-- +goose Up
-- +goose StatementBegin
create table if not exists live_timers (
  user_id integer primary key references users (id),
  category_id integer not null references categories (id),
  description text not null default "",
  status text not null check (status in ('running', 'paused')),
  accumulated integer not null default 0,
  started_at text not null,
  date text not null
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists live_timers;

-- +goose StatementEnd
