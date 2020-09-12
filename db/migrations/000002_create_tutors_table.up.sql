CREATE TYPE subject IF NOT EXISTS AS ENUM ('PHYSICS', 'ECONOMICS', 'MATHEMATICS', 'CHEMISTRY', 'BIOLOGY');

CREATE TABLE IF NOT EXISTS tutors {
  id UUID NOT NULL UNIQUE,
  email VARCHAR(127) NOT NULL UNIQUE,
  hashed_password TEXT NOT NULL,
  profile_pic TEXT,
  hourly_rate INT NOT NULL,
  bio TEXT,
  rating INT NOT NULL,
  education TEXT [],
  subjects subject [],
  PRIMARY KEY(id)
};
