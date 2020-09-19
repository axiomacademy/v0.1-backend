package db

import (
	"context"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Student struct {
	Id             string
	Username       string
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	ProfilePic     string
}

func (s *Student) Create(username string, firstName string, lastName string, email string, hashedPassword string, profile_pic string) error {

	// GENERATING UUID
	s.Id = "s:" + uuid.New()
	s.Username = username
	s.FirstName = firstName
	s.LastName = lastName
	s.Email = email
	s.HashedPassword = hashedPassword
	s.ProfilePic = profile_pic

	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO students (id, username, first_name, last_name, email, hashed_password, profile_pic) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.Username, s.FirstName, s.LastName, s.Email, s.HashedPassword, s.ProfilePic)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *Student) Update() error {
	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE students SET first_name = $2, last_name = $3, email = $4, hashed_password = $5, profile_pic = $6) WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.FirstName, s.LastName, s.Email, s.HashedPassword, s.ProfilePic)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (s *Student) GetById(id string) error {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic FROM students WHERE id = $1`

	if err := DbPool.QueryRow(context.Background(), sql, id).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic); err != nil {
		return err
	}

	return nil
}

func (s *Student) GetByUsername(username string) error {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic FROM students WHERE username = $1`

	if err := DbPool.QueryRow(context.Background(), sql, username).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic); err != nil {
		return err
	}

	return nil
}

func (s *Student) GetLessons() ([]Lesson, error) {
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

func (s *Student) ToModel() model.Student {
	return model.Student{ID: s.Id, Username: s.Username, FirstName: s.FirstName, LastName: s.LastName, Email: s.Email, ProfilePic: s.ProfilePic}
}
