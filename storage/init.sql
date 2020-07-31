CREATE TABLE IF NOT EXISTS users (
  id SERIAL NOT NULL PRIMARY KEY,
  username  text UNIQUE,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS chats (
  id SERIAL NOT NULL PRIMARY KEY,
  name  text UNIQUE,
  created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_chats (
  chat_id serial,
  user_id serial
);

CREATE TABLE IF NOT EXISTS messages (
  id SERIAL NOT NULL PRIMARY KEY,
  chatID serial,
  userID serial,
  text text,
  created_at timestamptz NOT NULL DEFAULT NOW()
);





