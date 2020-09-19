package db

import (
	"context"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

type Tutor struct {
	Id             string
	Email          string
	HashedPassword string
	ProfilePic     string
	HourlyRate     int
	Bio            string
	Rating         int
	Education      []string
	Subjects       []string
}

func (t *Tutor) Create(email string, password string, profile_pic string, hourly_rate int, rating int) error {

	// GENERATING UUID
	t.Id = "t:" + uuid.New()
	t.Email = email

	// Generating pasword hash
	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return err
	}

	t.HashedPassword = hashedPassword
	t.ProfilePic = profile_pic
	t.HourlyRate = hourly_rate
	t.Rating = rating

	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO tutors (id, email, hashed_password, profile_pic, hourly_rate, rating) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Rating)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (t *Tutor) Update() error {
	tx, err := DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE tutors SET email = $2, hashed_password = $3, profile_pic = $4, hourly_rate = $5, bio = $6, rating = $7, education = $8, subjects = $9 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Bio, t.Rating, t.Education, t.Subjects)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (t *Tutor) GetById(id string) error {
	sql := `SELECT id, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, subjects FROM tutors WHERE id = $1`

	if err := DbPool.QueryRow(context.Background(), sql, id).Scan(
		&t.Id,
		&t.Email,
		&t.HashedPassword,
		&t.ProfilePic,
		&t.HourlyRate,
		&t.Bio,
		&t.Rating,
		&t.Education,
		&t.Subjects); err != nil {
		return err
	}

	return nil
}

func (t *Tutor) GetLessons() ([]Lesson, error) {
	sql := `SELECT id, subject, tutor, student, duration, date, chat FROM lessons WHERE tutor = $1`

	var lessons []Lesson

	rows, err := DbPool.Query(context.Background(), sql, t.Id)
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
