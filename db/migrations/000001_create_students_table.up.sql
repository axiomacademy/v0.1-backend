CREATE TABLE IF NOT EXISTS students (
  id VARCHAR(38) NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email VARCHAR(127) NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  profile_pic TEXT NOT NULL DEFAULT '',
  push_token TEXT NOT NULL DEFAULT '',
  PRIMARY KEY(id)
);
