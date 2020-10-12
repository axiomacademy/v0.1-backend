package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

// Type mirror for Availability in the database
type Availability struct {
	Id        string
	Tutor     string
	StartTime time.Time
	EndTime   time.Time
}

// Creates an availability time block for a tutor.
// Takes a tutor id and recurring time period relative to epoch 0, ie. 0 epoch time corresponds to a recurring time at 12mid every Thursday
func (r *Repository) CreateAvailability(tid string, startTime time.Time, endTime time.Time) (Availability, error) {
	var a Availability

	a.Id = uuid.New()
	a.Tutor = tid
	a.StartTime = startTime
	a.EndTime = endTime

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return a, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO availabilities (id, tutor, start_time, end_time) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(context.Background(), sql, a.Id, a.Tutor, a.StartTime, a.EndTime)

	if err != nil {
		return a, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return a, err
	}

	return a, nil
}

// Deletes the availability for a tutor
// Takes a tutor id and recurring time period relative to epoch 0, ie. 0 epoch time corresponds to a recurring time at 12mid every Thursday
func (r *Repository) DeleteAvailability(tid string, startTime time.Time, endTime time.Time) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `DELETE FROM availabilities WHERE tutor = $1 AND start_time = $2 AND end_time = $3`
	_, err = tx.Exec(context.Background(), sql, tid, startTime, endTime)

	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Gets all tutor ids set to be potentially available within a timeslot sorted by affinity. Check the returned tutors with CheckTutorAvailability to verify with already booked lessons
// Takes a student id, subject id and a recurring time period relative to epoch 0 ie. 0 epoch time corresponds to a recurring time at 12mid every Thursday
func (r *Repository) GetAvailableTutors(sid string, subid string, startTime time.Time, endTime time.Time, limit int, offset int) ([]string, error) {
	sql := `
	SELECT availabilities.tutor
	FROM availabilities
	INNER JOIN affinity ON availabilities.tutor = affinity.tutor
	WHERE
		affinity.student = $1 AND
		affinity.subject = $2 AND
		availabilities.start_time <= timestamptz '$3' AND
		availabilities.end_time >= timestamptz '$4'
	ORDER_BY affinity.score DESC
	LIMIT $5 OFFSET $6`

	var tids []string
	rows, err := r.dbPool.Query(context.Background(), sql, sid, subid, startTime, endTime, limit, offset)
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

// Gets all tutor ids set to be potentially available within a timeslot shuffled randomly. Check the returned tutors with CheckTutorAvailability to verify with already booked lessons
// Takes a subject id,a recurring time period relative to epoch 0 ie. 0 epoch time corresponds to a recurring time at 12mid every Thursday, and a max count
func (r *Repository) GetRandomAvailableTutors(subid string, startTime time.Time, endTime time.Time, count int) ([]string, error) {
	sql := `
	SELECT availabilities.tutor
	FROM availabilities
	INNER JOIN subjects ON subjects.tutor = availabilities.tutor
	WHERE
		subjects.id = $1 AND
		availabilities.start_time <= timestamptz '$2' AND
		availabilities.end_time >= timestamptz '$3'
	TABLESAMPLE ($4 ROWS)
	LIMIT $5`

	var tids []string
	rows, err := r.dbPool.Query(context.Background(), sql, subid, startTime, endTime, count)
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

// Checks the lessons the tutor already has, to see if there are any availability clashes
func (r *Repository) CheckTutorAvailability(tid string, startTime time.Time, endTime time.Time) (bool, error) {
	sql := `SELECT id FROM lessons WHERE tutor = $1 AND scheduled = true AND start_time >= $2 AND start_time <= $3 AND end_time >= $2 AND end_time <= $3`

	var id string
	if err := r.dbPool.QueryRow(context.Background(), sql, tid, startTime, endTime).Scan(&id); err != nil {
		if err == pgx.ErrNoRows {
			return true, nil
		} else {
			return false, err
		}
	}

	return false, nil
}
