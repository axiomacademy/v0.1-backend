package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Lesson struct {
	Id       string
	Subject  Subject
	Summary  string
	Tutor    string
	Student  string
	Duration int
	Date     time.Time
	Chat     string
}

func (r *Repository) ToLessonModel(l Lesson) (model.Lesson, error) {

	s, err := r.GetStudentById(l.Student)
	if err != nil {
		return model.Lesson{}, err
	}

	t, err := r.GetTutorById(l.Tutor)
	if err != nil {
		return model.Lesson{}, err
	}

	rs := r.ToStudentModel(s)
	rt := r.ToTutorModel(t)
	rsub := l.Subject.ToSubjectModel()

	return model.Lesson{ID: l.Id, Subject: &rsub, Summary: l.Summary, Tutor: &rt, Student: &rs, Duration: l.Duration, Date: l.Date.String(), Chat: l.Chat}, nil
}

func (r *Repository) CreateLesson(subject Subject, tutor string, student string, duration int, date time.Time) (Lesson, error) {

	var l Lesson

	// GENERATING UUID
	l.Id = "l-" + uuid.New()
	l.Subject = subject
	l.Tutor = tutor
	l.Student = student
	l.Duration = duration
	l.Date = date

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return l, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO lessons (id, subject, tutor, student, duration, date) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject.Id, l.Tutor, l.Student, l.Duration, l.Date)

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

	sql := `UPDATE lessons SET subject = $2, summary = $3, tutor = $4, student = $5, duration = $6, date = $7, chat = $8 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject.Id, l.Summary, l.Tutor, l.Student, l.Duration, l.Date, l.Chat)

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
	var date pgtype.Timestamptz
	var l Lesson

	sql := `SELECT lessons.id, subjects.id, subjects.name, subjects.standard, lessons.summary, lessons.tutor, lessons.student, lessons.duration, lessons.date, lessons.chat FROM lessons INNER JOIN subjects ON lessons.subject = subjects.id WHERE lessons.id = $1`
	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(&l.Id, &l.Subject.Id, &l.Subject.Name, &l.Subject.Standard, &l.Summary, &l.Tutor, &l.Student, &l.Duration, &date, &l.Chat); err != nil {
		return l, err
	}

	date.AssignTo(&l.Date)
	return l, nil
}
