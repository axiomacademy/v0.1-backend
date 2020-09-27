CREATE TABLE IF NOT EXISTS affinity (
  tutor VARCHAR(38) NOT NULL,
  student VARCHAR(38) NOT NULL,
  subject subject NOT NULL,
  score INT NOT NULL DEFAULT 0,
  PRIMARY KEY(tutor, student, subject),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id)
);
