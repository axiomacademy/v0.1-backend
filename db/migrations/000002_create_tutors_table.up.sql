CREATE TYPE subject AS ENUM ('PHYSICS', 'ECONOMICS', 'MATHEMATICS', 'CHEMISTRY', 'BIOLOGY');

CREATE TYPE status AS ENUM ('AVAILABLE', 'BUSY')

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
  subjects subject [],
	status status,
	last_seen TIMESTAMP,
  PRIMARY KEY(id)
);
