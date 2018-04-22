-- +migrate Up

CREATE TABLE users
(
  id   serial PRIMARY KEY NOT NULL,
  name text               NOT NULL
);
CREATE INDEX users_name_index
  ON users (name);

-- +migrate Down
drop table users;
