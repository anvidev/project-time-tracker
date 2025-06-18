-- +goose Up
-- +goose StatementBegin
create table if not exists users_hours (
  user_id integer not null,
  weekday integer not null,
  hours text not null,
  primary key (user_id, weekday),
  foreign key (user_id) references users (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
drop table if exists users_hours;

-- +goose StatementEnd
