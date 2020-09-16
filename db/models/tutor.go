package models

import (
	"context"
	db "github.com/solderneer/axiom-backend/db"
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

func (t Tutor) Create() error {
	tx, err := db.DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO tutor (id, email, hashed_password, profile_pic, hourly_rate, rating) VALUES ($1, $2, $3, $4, $5, $6)`
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

func (t Tutor) Update() error {
	tx, err := db.DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE tutor SET email = $2, hashed_password = $3, profile_pic = $4, hourly_rate = $5, bio = $6, rating = $7, education = $8, subjects = $9 WHERE id = $1`
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

func (t Tutor) GetById(id string) error {
	sql := `SELECT id, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, subjects FROM student WHERE id = $1`

	if err := db.DbPool.QueryRow(context.Background(), sql, id).Scan(
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
