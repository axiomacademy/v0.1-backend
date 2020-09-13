CREATE TABLE IF NOT EXISTS lessons {
  id UUID NOT NULL UNIQUE,
  subject subject NOT NULL,
  tutor UUID NOT NULL,
  student UUID NOT NULL,
  duration INT NOT NULL,
  date TIMESTAMP NOT NULL,
  chat TEXT,
  PRIMARY KEY(id),
  CONSTRAINT fk_tutor
    FOREIGN KEY(tutor)
      REFERENCES tutors(id),
  CONSTRAINT fk_student
    FOREIGN KEY(student)
      REFERENCES students(id)
};
  
