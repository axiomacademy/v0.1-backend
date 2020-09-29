package db

import (
	"context"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Tutor struct {
	Id             string
	Username       string
	FirstName      string
	LastName       string
	Email          string
	HashedPassword string
	ProfilePic     string
	HourlyRate     int
	Bio            string
	Rating         int
	Education      []string
	Subject        Subject
}

func (r *Repository) CreateTutor(username string, firstName string, lastName string, email string, hashedPassword string, profile_pic string, hourly_rate int, rating int, bio string, education []string, subject Subject) (Tutor, error) {

	var t Tutor

	// GENERATING UUID
	t.Id = "t:" + uuid.New()
	t.Username = username
	t.FirstName = firstName
	t.LastName = lastName
	t.Email = email
	t.HashedPassword = hashedPassword
	t.ProfilePic = profile_pic
	t.HourlyRate = hourly_rate
	t.Rating = rating
	t.Bio = bio
	t.Education = education
	t.Subject = subject

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return t, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO tutors (id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, rating, bio, education, subject) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.Username, t.FirstName, t.LastName, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Rating, t.Bio, t.Education, t.Subject)

	if err != nil {
		return t, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return t, err
	}

	return t, nil
}

func (r *Repository) UpdateTutor(t Tutor) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE tutors SET first_name = $2, last_name = $3, email = $4, hashed_password = $5, profile_pic = $6, hourly_rate = $7, bio = $8, rating = $9, education = $10, subject = $11  WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.FirstName, t.LastName, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Bio, t.Rating, t.Education, t.Subject)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetTutorById(id string) (Tutor, error) {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, subject FROM tutors WHERE id = $1`

	var t Tutor

	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(
		&t.Id,
		&t.Username,
		&t.FirstName,
		&t.LastName,
		&t.Email,
		&t.HashedPassword,
		&t.ProfilePic,
		&t.HourlyRate,
		&t.Bio,
		&t.Rating,
		&t.Education,
		&t.Subject); err != nil {
		return t, err
	}

	return t, nil
}

func (r *Repository) GetTutorByUsername(username string) (Tutor, error) {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, subject FROM tutors WHERE username = $1`

	var t Tutor

	if err := r.dbPool.QueryRow(context.Background(), sql, username).Scan(
		&t.Id,
		&t.Username,
		&t.FirstName,
		&t.LastName,
		&t.Email,
		&t.HashedPassword,
		&t.ProfilePic,
		&t.HourlyRate,
		&t.Bio,
		&t.Rating,
		&t.Education,
		&t.Subject); err != nil {
		return t, err
	}

	return t, nil
}

func (r *Repository) GetTutorLessons(tid string) ([]Lesson, error) {
	sql := `SELECT id, subject, summary, tutor, student, duration, date, chat FROM lessons WHERE tutor = $1`

	var lessons []Lesson

	rows, err := r.dbPool.Query(context.Background(), sql, tid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var lesson Lesson
		err := rows.Scan(&lesson.Id, &lesson.Subject, &lesson.Summary, &lesson.Tutor, &lesson.Student, &lesson.Duration, &lesson.Date, &lesson.Chat)

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

func (r *Repository) ToTutorModel(t Tutor) model.Tutor {
	subject := model.Subject{Name: model.SubjectName(t.Subject.Name), Level: model.SubjectLevel(t.Subject.Level)}

	return model.Tutor{ID: t.Id, Username: t.Username, FirstName: t.FirstName, LastName: t.LastName, Email: t.Email, ProfilePic: t.ProfilePic, HourlyRate: t.HourlyRate, Bio: t.Bio, Rating: t.Rating, Education: t.Education, Subject: &subject}
}
