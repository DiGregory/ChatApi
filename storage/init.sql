CREATE TABLE IF NOT EXISTS users (
  id         SERIAL      NOT NULL PRIMARY KEY,
  username   text UNIQUE,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chats (
  id         SERIAL      NOT NULL PRIMARY KEY,
  name       text UNIQUE,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_chats (
  chat_id integer,
  user_id integer
);

CREATE TABLE IF NOT EXISTS messages (
  id         SERIAL      NOT NULL PRIMARY KEY,
  chat_id    integer,
  user_id    integer,
  text       text,
  created_at timestamptz NOT NULL DEFAULT NOW()
);





