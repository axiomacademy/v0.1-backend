CREATE TABLE IF NOT EXISTS lessons (
  id VARCHAR(38) NOT NULL UNIQUE,
  subject VARCHAR(38) NOT NULL,
  summary TEXT,
  tutor VARCHAR(38) NOT NULL,
  student VARCHAR(38) NOT NULL,
  duration INT NOT NULL,
  date TIMESTAMP NOT NULL,
  chat TEXT,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id),
  CONSTRAINT fk_subject
    FOREIGN KEY(subject)
      REFERENCES subjects(id)
);
  
