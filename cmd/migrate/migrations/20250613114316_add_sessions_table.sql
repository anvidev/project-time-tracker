-- +goose Up
-- +goose StatementBegin
create table if not exists sessions (
  user_id integer not null references users (id),
  expires_at text not null
) random rowid;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists sessions;

-- +goose StatementEnd
