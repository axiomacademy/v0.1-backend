CREATE TABLE IF NOT EXISTS students (
  id VARCHAR(38) NOT NULL UNIQUE,
  email VARCHAR(127) NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  profile_pic TEXT,
  PRIMARY KEY(id)
);
