package db

import (
	"context"
	"time"

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

// Convert from db.Student to model.Student
func (r *Repository) ToStudentModel(s Student) model.Student {
	return model.Student{ID: s.Id, Username: s.Username, FirstName: s.FirstName, LastName: s.LastName, Email: s.Email, ProfilePic: s.ProfilePic}
}

// Creates a new student, returns a db.Student
func (r *Repository) CreateStudent(username string, firstName string, lastName string, email string, hashedPassword string, profilePic string) (Student, error) {

	var s Student

	// GENERATING UUID
	s.Id = "s:" + uuid.New()
	s.Username = username
	s.FirstName = firstName
	s.LastName = lastName
	s.Email = email
	s.HashedPassword = hashedPassword
	s.ProfilePic = profilePic

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

// Update student, takes an existing student. Only can update firstname, lastname, email, profile picture and APN push token
func (r *Repository) UpdateStudent(s Student) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE students SET first_name = $2, last_name = $3, email = $4, profile_pic = $5, push_token = $6 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, s.Id, s.FirstName, s.LastName, s.Email, s.ProfilePic, s.PushToken)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Gets the student by student UUID
func (r *Repository) GetStudentById(id string) (Student, error) {

	var s Student

	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, push_token FROM students WHERE id = $1`

	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic, &s.PushToken); err != nil {
		return s, err
	}

	return s, nil
}

// Gets student by student username
func (r *Repository) GetStudentByUsername(username string) (Student, error) {

	var s Student

	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, push_token FROM students WHERE username = $1`

	if err := r.dbPool.QueryRow(context.Background(), sql, username).Scan(&s.Id, &s.Username, &s.FirstName, &s.LastName, &s.Email, &s.HashedPassword, &s.ProfilePic, &s.PushToken); err != nil {
		return s, err
	}

	return s, nil
}

// Gets all student lessons, paginated by startTime and endTime
func (r *Repository) GetStudentLessons(sid string, startTime time.Time, endTime time.Time) ([]Lesson, error) {
	sql := `SELECT id, subject, tutor, student, scheduled, period FROM lessons WHERE student = $1 and $2 @> period`

	var lessons []Lesson

	period := getTstzrange(startTime, endTime)

	rows, err := r.dbPool.Query(context.Background(), sql, sid, period)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var lesson Lesson
		var period pgtype.Tstzrange
		var sid string

		err := rows.Scan(&lesson.Id, &sid, &lesson.Tutor, &lesson.Student, &lesson.Scheduled, &period)

		if err != nil {
			return nil, err
		}

		period.Upper.AssignTo(&lesson.StartTime)
		period.Lower.AssignTo(&lesson.EndTime)

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

func (r *Repository) IsStudentInLesson(sid string, lid string) (bool, error) {
	sql := `SELECT 1 FROM lessons WHERE id = $1 AND student = $2`

	rows, err := r.dbPool.Query(context.Background(), sql, lid, sid)
	if err != nil {
		return false, err
	}

	return rows.Next(), nil
}
