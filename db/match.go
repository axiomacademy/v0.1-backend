package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Match struct {
	Id        string
	Status    string
	Scheduled bool
	Tutor     string
	Student   string
	Subject   string
	StartTime time.Time
	EndTime   time.Time
	Lesson    string
}

// Converts a db.Match to model.Match
func (r *Repository) ToMatchModel(m Match) (model.Match, error) {
	s, err := r.GetStudentById(m.Student)
	if err != nil {
		return model.Match{}, err
	}

	t, err := r.GetTutorById(m.Tutor)
	if err != nil {
		return model.Match{}, err
	}

	sub, err := r.GetSubjectById(m.Subject)
	if err != nil {
		return model.Match{}, err
	}

	rs := r.ToStudentModel(s)
	rt, err := r.ToTutorModel(t)
	if err != nil {
		return model.Match{}, err
	}
	rsub := r.ToSubjectModel(sub)

	return model.Match{ID: m.Id, Status: m.Status, Scheduled: m.Scheduled, Tutor: &rt, Student: &rs, Subject: &rsub, StartTime: &m.StartTime, EndTime: &m.EndTime}, nil
}

// Create a new on-demand match. Takes in the status string, student UUID string, subject UUID string
func (r *Repository) CreateOnDemandMatch(status string, sid string, subid string) (Match, error) {
	var m Match
	m.Id = uuid.New()
	m.Status = status
	m.Scheduled = false
	m.Student = sid
	m.Subject = subid

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return m, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO matchings (id, status, scheduled, student, subject) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(context.Background(), sql, m.Id, m.Status, m.Scheduled, m.Student, m.Subject)

	if err != nil {
		return m, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return m, err
	}

	return m, nil
}

// Create a new scheduled match process
// Takes a status string, student UUID string, tutor UUID string, subject UUID string, startTime and endTime in absolute time.Time
func (r *Repository) CreateScheduledMatch(status string, sid string, tid string, subid string, startTime time.Time, endTime time.Time) (Match, error) {
	var m Match
	m.Id = uuid.New()
	m.Status = status
	m.Scheduled = false
	m.Student = sid
	m.Tutor = tid
	m.Subject = subid
	m.StartTime = startTime
	m.EndTime = endTime

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return m, err
	}

	defer tx.Rollback(context.Background())

	period := getTstzrange(startTime, endTime)

	sql := `INSERT INTO matchings (id, status, scheduled, student, tutor, subject, period) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err = tx.Exec(context.Background(), sql, m.Id, m.Status, m.Scheduled, m.Student, m.Tutor, m.Subject, period)

	if err != nil {
		return m, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return m, err
	}

	return m, nil
}

// Updates match row, only allows update of status and lesson column
// Takes in the updated match struct and then updates database state to match
func (r *Repository) UpdateMatch(m Match) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE matchings SET status = $2, lesson $3 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, m.Id, m.Status, m.Lesson)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Gets the match struct based on the match UUID
func (r *Repository) GetMatchById(mid string) (Match, error) {
	sql := `SELECT id, status, scheduled, tutor, student, subject, period, lesson FROM matchings WHERE id = $1`
	var period pgtype.Tstzrange
	var m Match

	if err := r.dbPool.QueryRow(context.Background(), sql, mid).Scan(&m.Id, &m.Status, &m.Scheduled, &m.Tutor, &m.Student, &m.Subject, &period, &m.Lesson); err != nil {
		return m, err
	}

	period.Upper.AssignTo(&m.EndTime)
	period.Lower.AssignTo(&m.StartTime)

	return m, nil
}

func (r *Repository) GetTutorPendingMatches(tid string) ([]Match, error) {
	sql := `SELECT id, status, scheduled, tutor, student, subject, period, lesson FROM matchings WHERE tutor = $1 AND status = $2`

	var matches []Match

	rows, err := r.dbPool.Query(context.Background(), sql, tid, "MATCHING")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var m Match
		var period pgtype.Tstzrange

		err := rows.Scan(&m.Id, &m.Status, &m.Scheduled, &m.Tutor, &m.Student, &m.Subject, &period, &m.Lesson)

		if err != nil {
			return nil, err
		}

		period.Upper.AssignTo(&m.EndTime)
		period.Lower.AssignTo(&m.StartTime)

		matches = append(matches, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return matches, nil
}

// Gets all the pending matches for a student, takes in a student UUID
func (r *Repository) GetStudentPendingMatches(sid string) ([]Match, error) {
	sql := `SELECT id, status, scheduled, tutor, student, subject, period, lesson FROM matchings WHERE student = $1 AND status = $2`

	var matches []Match

	rows, err := r.dbPool.Query(context.Background(), sql, sid, "MATCHING")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var m Match
		var period pgtype.Tstzrange

		err := rows.Scan(&m.Id, &m.Status, &m.Scheduled, &m.Tutor, &m.Student, &m.Subject, &period, &m.Lesson)

		if err != nil {
			return nil, err
		}

		period.Upper.AssignTo(&m.EndTime)
		period.Lower.AssignTo(&m.StartTime)

		matches = append(matches, m)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return matches, nil
}
