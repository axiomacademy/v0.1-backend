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

// Creates a new affinity, takes in the tutor UUID, student UUID and subject UUID
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

// Deletes an existing affinity, takes in an existing affinity UUID
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

// Updates an existing affinity, takes in an existing retrieved affinity
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

// Get the affinity given the tutor UUID, string UUID and subject UUID
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

// Get random on-demand matches given a subject based on online status. Limited by count
func (r *Repository) GetOnlineRandomMatches(subid string, count int) ([]string, error) {
	sql := `
	SELECT tutors.id 
	FROM teaching
	INNER JOIN tutors ON tutors.id = teaching.tutor
	WHERE
		tutors.last_seen > $1 AND
		tutors.status = AVAILABLE AND
		teaching.subject = $2
	TABLESAMPLE ($3 ROWS)
	LIMIT $4
	`

	var tids []string

	exp, err := time.Now().Add(time.Minute * 1).MarshalText()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(context.Background(), sql, exp, subid, count*100, count)
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

// Get online on-demand matches, based and sorted by affinity. Limited by the count parameter
func (r *Repository) GetOnlineAffinityMatches(sid string, subid string, count int) ([]string, error) {
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
	LIMIT $4`

	var tids []string

	exp, err := time.Now().Add(time.Minute * 1).MarshalText()
	if err != nil {
		return nil, err
	}

	rows, err := r.dbPool.Query(context.Background(), sql, sid, subid, exp, count)
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
