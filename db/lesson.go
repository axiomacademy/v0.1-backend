package db

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Lesson struct {
	Id       string
	Subject  string
	Summary  string
	Tutor    string
	Student  string
	Duration int
	Date     time.Time
	Chat     string
}

func (l *Lesson) Create(subject string, tutor string, student string, duration int, date time.Time) error {
	// GENERATING UUID
	l.Id = "l-" + uuid.New()
	l.Subject = subject
	l.Tutor = tutor
	l.Student = student
	l.Duration = duration
	l.Date = date

	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO lessons (id, subject, tutor, student, duration, date) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject, l.Tutor, l.Student, l.Duration, l.Date)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (l *Lesson) Update() error {
	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE lessons SET subject = $2, summary = $3, tutor = $4, student = $5, duration = $6, date = $7, chat = $8 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject, l.Summary, l.Tutor, l.Student, l.Duration, l.Date, l.Chat)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (l *Lesson) GetById(id string) error {
	sql := `SELECT id, subject, summary, tutor, student, duration, date, chat FROM lessons WHERE id = $1`

	if err := DbPool.QueryRow(context.Background(), sql, id).Scan(&l.Id, &l.Subject, &l.Summary, &l.Tutor, &l.Student, &l.Duration, &l.Date, &l.Chat); err != nil {
		return err
	}

	return nil
}

func (l *Lesson) ToModel() model.Lesson {
	s := Student{}
	t := Tutor{}

	s.GetById(l.Student)
	t.GetById(l.Tutor)

	rs := s.ToModel()
	rt := t.ToModel()

	return model.Lesson{ID: l.Id, Subject: l.Subject, Summary: l.Summary, Tutor: &rt, Student: &rs, Duration: l.Duration, Date: l.Date.String(), Chat: l.Chat}
}
