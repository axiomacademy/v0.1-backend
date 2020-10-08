CREATE TABLE IF NOT EXISTS notifications (
  id VARCHAR(36),
  tutor VARCHAR(38),
  student VARCHAR(38),
  image TEXT,
  title TEXT NOT NULL,
  subtitle TEXT NOT NULL,
  read BOOL NOT NULL DEFAULT FALSE,
  created TIMESTAMPTZ NOT NULL DEFAULT,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id)
);

