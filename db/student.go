package db

import (
	"context"
	"github.com/pborman/uuid"
)

type Student struct {
	Id             string
	Email          string
	HashedPassword string
	ProfilePic     string
}

func (s *Student) Create(email string, hashed_password string, profile_pic string) error {

	// GENERATING UUID
	s.Id = "s-" + uuid.New()
	s.Email = email
	s.HashedPassword = hashed_password
	s.ProfilePic = profile_pic

	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO students (id, email, hashed_password, profile_pic) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.Email, s.HashedPassword, s.ProfilePic)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s Student) Update() error {
	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE students SET email = $2, hashed_password = $3, profile_pic = $4) WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.Email, s.HashedPassword, s.ProfilePic)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s Student) GetById(id string) error {
	sql := `SELECT id, email, hashed_password, profile_pic FROM students WHERE id = $1`

	if err := DbPool.QueryRow(context.Background(), sql, id).Scan(&s.Id, &s.Email, &s.HashedPassword, &s.ProfilePic); err != nil {
		return err
	}

	return nil
}

func (s Student) GetLessons() ([]Lesson, error) {
	sql := `SELECT id, subject, tutor, student, duration, date, chat FROM lessons WHERE student = $1`

	var lessons []Lesson

	rows, err := DbPool.Query(context.Background(), sql, s.Id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.Id, &lesson.Subject, &lesson.Tutor, &lesson.Student, &lesson.Duration, &lesson.Date, &lesson.Chat)

		if err != nil {
			return nil, err
		}

		lessons = append(lessons, lesson)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return lessons, nil
}
