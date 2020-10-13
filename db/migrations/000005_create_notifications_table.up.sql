CREATE TABLE IF NOT EXISTS notifications (
  id VARCHAR(38),
  tutor VARCHAR(38),
  student VARCHAR(38),
  title TEXT NOT NULL,
  subtitle TEXT NOT NULL,
  image TEXT,
  read BOOL NOT NULL DEFAULT FALSE,
  created TIMESTAMPTZ NOT NULL,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id)
);

CREATE TABLE IF NOT EXISTS matchings (
  id VARCHAR(38),
  status TEXT NOT NULL,
  scheduled BOOLEAN NOT NULL,
  tutor VARCHAR(38),
  student VARCHAR(38),
  subject VARCHAR(38) NOT NULL,
  period TSTZRANGE NOT NULL,
  lesson VARCHAR(38),
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id),
  CONSTRAINT fk_subject
    FOREIGN KEY(subject)
      REFERENCES subjects(id),
  CONSTRAINT fk_lesson
    FOREIGN KEY(lesson)
      REFERENCES lessons(id)
);


