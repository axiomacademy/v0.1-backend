package db

import (
	"context"
	"time"
)

type Affinity struct {
	Tutor   string
	Student string
	Subject string
	Score   int
}

func (r *Repository) CreateAffinity(tid string, sid string, subid string) (Affinity, error) {
	var a Affinity
	a.Tutor = tid
	a.Student = sid
	a.Subject = subid
	a.Score = 0

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return a, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO affinity (tutor, student, subject, score) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(context.Background(), sql, a.Tutor, a.Student, a.Subject, a.Score)

	if err != nil {
		return a, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return a, err
	}

	return a, nil
}

func (r *Repository) DeleteAffinity(a Affinity) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `DELETE FROM affinity WHERE tutor = $1 AND student = $2 AND subject = $3`
	_, err = tx.Exec(context.Background(), sql, a.Tutor, a.Student, a.Subject)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateAffinity(a Affinity) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE affinity SET score = $4 WHERE tutor = $1 AND student = $2 AND subject = $3`
	_, err = tx.Exec(context.Background(), sql, a.Tutor, a.Student, a.Subject, a.Score)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAffinity(tid string, sid string, subid string) (Affinity, error) {
	var a Affinity

	a.Tutor = tid
	a.Student = sid
	a.Subject = subid

	sql := `SELECT score FROM affinity WHERE tutor = $1, student = $2, subject = $3`
	if err := r.dbPool.QueryRow(context.Background(), sql, a.Tutor, a.Student, a.Subject).Scan(&a.Score); err != nil {
		return a, err
	}

	return a, nil
}

func (r *Repository) GetOnlineRandomMatches(subject Subject, count int) ([]string, error) {
	sql := `
	SELECT tutors.id 
	FROM teaching
	INNER JOIN tutors ON tutors.id = teaching.tutor
	WHERE
		tutors.last_seen > timestamptz '$1' AND
		tutors.status = AVAILABLE
	TABLESAMPLE ($2 ROWS)
	LIMIT $3
	`

	var tids []string

	exp, err := time.Now().Add(time.Minute * 1).MarshalText()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(context.Background(), sql, exp, 100, count)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var tid string
		if err := rows.Scan(&tid); err != nil {
			return nil, err
		}

		tids = append(tids, tid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tids, nil
}

func (r *Repository) GetOnlineAffinityMatches(sid string, subject Subject) ([]string, error) {
	sql := `
	SELECT affinity.tutor 
	FROM affinity 
	INNER JOIN tutors ON tutors.id = affinity.tutor 
	WHERE 
		affinity.student = $1 AND
		affinity.subject = $2 AND
		tutors.last_seen > timestamptz '$3' AND
		tutors.status = AVAILABLE
	ORDER_BY affinity.score DESC
	LIMIT $4 OFFSET $5`

	var tids []string

	exp, err := time.Now().Add(time.Minute * 1).MarshalText()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(context.Background(), sql, sid, subject.Id, exp, 10, 0)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var tid string
		if err := rows.Scan(&tid); err != nil {
			return nil, err
		}

		tids = append(tids, tid)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tids, nil
}
