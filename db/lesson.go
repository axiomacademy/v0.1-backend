package db

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Lesson struct {
	Id           string
	Subject      string
	SubjectLevel string
	Summary      string
	Tutor        string
	Student      string
	Duration     int
	Date         time.Time
	Chat         string
}

func (r *Repository) ToLessonModel(l Lesson) model.Lesson {
	s := Student{}
	t := Tutor{}

	r.GetStudentById(l.Student)
	r.GetTutorById(l.Tutor)

	rs := r.ToStudentModel(s)
	rt := r.ToTutorModel(t)

	return model.Lesson{ID: l.Id, Subject: l.Subject, SubjectLevel: l.SubjectLevel, Summary: l.Summary, Tutor: &rt, Student: &rs, Duration: l.Duration, Date: l.Date.String(), Chat: l.Chat}
}

func (r *Repository) CreateLesson(subject string, subject_level string, tutor string, student string, duration int, date time.Time) (Lesson, error) {

	var l Lesson

	// GENERATING UUID
	l.Id = "l-" + uuid.New()
	l.Subject = subject
	l.SubjectLevel = subject_level
	l.Tutor = tutor
	l.Student = student
	l.Duration = duration
	l.Date = date

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return l, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO lessons (id, subject, subject_level, tutor, student, duration, date) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject, l.SubjectLevel, l.Tutor, l.Student, l.Duration, l.Date)

	if err != nil {
		return l, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return l, err
	}

	return l, nil
}

func (r *Repository) UpdateLesson(l Lesson) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE lessons SET subject = $2, subject_level, $3, summary = $4, tutor = $5, student = $6, duration = $7, date = $8, chat = $9 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject, l.SubjectLevel, l.Summary, l.Tutor, l.Student, l.Duration, l.Date, l.Chat)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLessonById(id string) (Lesson, error) {
	var l Lesson
	sql := `SELECT id, subject, subject_level, summary, tutor, student, duration, date, chat FROM lessons WHERE id = $1`

	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(&l.Id, &l.Subject, &l.SubjectLevel, &l.Summary, &l.Tutor, &l.Student, &l.Duration, &l.Date, &l.Chat); err != nil {
		return l, err
	}

	return l, nil
}
