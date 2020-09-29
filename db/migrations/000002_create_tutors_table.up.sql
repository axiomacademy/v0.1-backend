CREATE TYPE subject AS (
  name TEXT,
  'level' TEXT
);

CREATE TABLE IF NOT EXISTS tutors (
  id VARCHAR(38) NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email VARCHAR(127) NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  profile_pic TEXT,
  hourly_rate INT NOT NULL,
  bio TEXT,
  rating INT NOT NULL,
  education TEXT [],
  subject subject,
  status VARCHAR(10) NOT NULL DEFAULT 'OFFLINE',
  last_seen TIMESTAMP,
  PRIMARY KEY(id)
);
