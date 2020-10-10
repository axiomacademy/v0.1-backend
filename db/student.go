package db

import (
	"context"

	"github.com/jackc/pgtype"
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
	PushToken      string
}

func (r *Repository) ToStudentModel(s Student) model.Student {
	return model.Student{ID: s.Id, Username: s.Username, FirstName: s.FirstName, LastName: s.LastName, Email: s.Email, ProfilePic: s.ProfilePic}
}

func (r *Repository) CreateStudent(username string, firstName string, lastName string, email string, hashedPassword string, profile_pic string) (Student, error) {

	var s Student

	// GENERATING UUID
	s.Id = "s:" + uuid.New()
	s.Username = username
	s.FirstName = firstName
	s.LastName = lastName
	s.Email = email
	s.HashedPassword = hashedPassword
	s.ProfilePic = profile_pic

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return s, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO students (id, username, first_name, last_name, email, hashed_password, profile_pic) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.Username, s.FirstName, s.LastName, s.Email, s.HashedPassword, s.ProfilePic)

	if err != nil {
		return s, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return s, err
	}

	return s, nil
}

func (r *Repository) UpdateStudent(s Student) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE students SET first_name = $2, last_name = $3, email = $4, hashed_password = $5, profile_pic = $6, push_token = $7 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.FirstName, s.LastName, s.Email, s.HashedPassword, s.ProfilePic, s.PushToken)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetStudentById(id string) (Student, error) {

	var s Student

	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, push_token FROM students WHERE id = $1`

	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic, &s.PushToken); err != nil {
		return s, err
	}

	return s, nil
}

func (r *Repository) GetStudentByUsername(username string) (Student, error) {

	var s Student

	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, push_token FROM students WHERE username = $1`

	if err := r.dbPool.QueryRow(context.Background(), sql, username).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic, &s.PushToken); err != nil {
		return s, err
	}

	return s, nil
}

func (r *Repository) GetStudentLessons(sid string) ([]Lesson, error) {
	sql := `SELECT id, subject, tutor, student, duration, date, chat FROM lessons WHERE student = $1`

	var lessons []Lesson

	rows, err := r.dbPool.Query(context.Background(), sql, sid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var lesson Lesson
		var date pgtype.Timestamptz
		var sid string

		err := rows.Scan(&lesson.Id, &sid, &lesson.Tutor, &lesson.Student, &lesson.Duration, &date, &lesson.Chat)

		if err != nil {
			return nil, err
		}

		date.AssignTo(&lesson.Date)
		if lesson.Subject, err = r.GetSubjectById(sid); err != nil {
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
