CREATE TABLE IF NOT EXISTS tutors (
  id VARCHAR(38) NOT NULL UNIQUE,
  username TEXT NOT NULL UNIQUE,
  first_name TEXT NOT NULL,
  last_name TEXT NOT NULL,
  email VARCHAR(127) NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  profile_pic TEXT NOT NULL DEFAULT '',
  hourly_rate INT NOT NULL,
  bio TEXT NOT NULL DEFAULT '',
  rating INT NOT NULL,
  education TEXT [] NOT NULL DEFAULT {''},
  status VARCHAR(12) NOT NULL DEFAULT 'UNAVAILABLE',
	last_seen TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  push_token TEXT NOT NULL DEFAULT '',
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

CREATE TABLE IF NOT EXISTS availabilities (
  id VARCHAR(38) NOT NULL UNIQUE,
  tutor VARCHAR(38) NOT NULL,
  period TSTZRANGE NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id)
);
