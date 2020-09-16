package models

import (
	"context"
	db "github.com/solderneer/axiom-backend/db"
)

type Student struct {
	Id             string
	Email          string
	HashedPassword string
	ProfilePic     string
}

func (s Student) Create() error {
	tx, err := db.DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO student (id, email, hashed_password, profile_pic) VALUES ($1, $2, $3, $4)`
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
	tx, err := db.DbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE student SET email = $2, hashed_password = $3, profile_pic = $4) WHERE id = $1`
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
	sql := `SELECT id, email, hashed_password, profile_pic FROM student WHERE id = $1`

	if err := db.DbPool.QueryRow(context.Background(), sql, id).Scan(&s.Id, &s.Email, &s.HashedPassword, &s.ProfilePic); err != nil {
		return err
	}

	return nil
}
