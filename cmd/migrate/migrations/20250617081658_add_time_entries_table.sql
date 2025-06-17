-- +goose Up
-- +goose StatementBegin
create table if not exists time_entries (
  id integer primary key,
  category_id integer not null references categories (id),
  user_id integer not null references users (id),
  date text not null,
  duration text not null,
  description text not null default ""
);

create index if not exists idx_time_entries_user_id on time_entries (user_id);

create index if not exists idx_time_entries_category_id on time_entries (category_id);

create index if not exists idx_time_entries_date on time_entries (date);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop index if exists idx_time_entries_user_id;

drop index if exists idx_time_entries_category_id;

drop index if exists idx_time_entries_date;

drop table if exists time_entries;

-- +goose StatementEnd
