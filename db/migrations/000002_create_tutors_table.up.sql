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
  status VARCHAR(12),
	last_seen TIMESTAMPTZ,
  push_token TEXT,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS subjects (
  id VARCHAR(38) NOT NULL UNIQUE,
  name TEXT NOT NULL,
  standard TEXT NOT NULL,
  PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS teaching (
  tutor VARCHAR(38) NOT NULL,
  subject VARCHAR(38) NOT NULL,
  PRIMARY KEY(tutor, subject),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_subject
    FOREIGN KEY(subject)
      REFERENCES subjects(id)
);

CREATE TABLE IF NOT EXISTS timeblocks (
  id VARCHAR(38) NOT NULL UNIQUE,
  tutor VARCHAR(38) NOT NULL,
  start_time TIMESTAMPTZ NOT NULL,
  end_time TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id)
);
