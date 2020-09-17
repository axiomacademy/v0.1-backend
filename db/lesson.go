package db

import (
	"context"
	"time"
)

type Lesson struct {
	Id       string
	Subject  string
	Tutor    string
	Student  string
	Duration int
	Date     time.Time
	Chat     string
}

func (l Lesson) Create() error {
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

func (l Lesson) Update() error {
	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE lessons SET subject = $2, tutor = $3, student = $4, duration = $5, date = $6, chat = $7 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject, l.Tutor, l.Student, l.Duration, l.Date, l.Chat)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (l Lesson) GetById(id string) error {
	sql := `SELECT id, subject, tutor, student, duration, date, chat FROM lessons WHERE id = $1`

	if err := DbPool.QueryRow(context.Background(), sql, id).Scan(&l.Id, &l.Subject, &l.Tutor, &l.Student, &l.Duration, &l.Date, &l.Chat); err != nil {
		return err
	}

	return nil
}
