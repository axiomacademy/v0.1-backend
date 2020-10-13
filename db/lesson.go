package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Lesson struct {
	Id        string
	Subject   Subject
	Summary   string
	Tutor     string
	Student   string
	Scheduled bool
	StartTime time.Time
	EndTime   time.Time
}

func (r *Repository) ToLessonModel(l Lesson) (model.Lesson, error) {

	s, err := r.GetStudentById(l.Student)
	if err != nil {
		return model.Lesson{}, err
	}

	t, err := r.GetTutorById(l.Tutor)
	if err != nil {
		return model.Lesson{}, err
	}

	rs := r.ToStudentModel(s)
	rt := r.ToTutorModel(t)
	rsub := r.ToSubjectModel(l.Subject)

	return model.Lesson{ID: l.Id, Subject: &rsub, Summary: l.Summary, Tutor: &rt, Student: &rs, Scheduled: l.Scheduled, StartTime: l.StartTime, EndTime: l.EndTime}, nil
}

// Creates a new lesson in the database
// Accepts a subject, tutor ID, student ID, scheduled status, and a startTime + endTime
// For a scheduled lesson, startTime/endTime are releative to 0 epoch time and represent points in a week. ie. 0 epoch time is Thursday 12 midnight
func (r *Repository) CreateLesson(subject Subject, tutor string, student string, scheduled bool, startTime time.Time, endTime time.Time) (Lesson, error) {

	var l Lesson

	// GENERATING UUID
	l.Id = "l-" + uuid.New()
	l.Subject = subject
	l.Tutor = tutor
	l.Student = student
	l.Scheduled = scheduled
	l.StartTime = startTime
	l.EndTime = endTime

	period := getTstzrange(l.StartTime, l.EndTime)

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return l, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO lessons (id, subject, tutor, student, scheduled, period) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject.Id, l.Tutor, l.Student, l.Scheduled, period)

	if err != nil {
		return l, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return l, err
	}

	return l, nil
}

// Updates an existing lesson in the database
// Accepts a lesson object
func (r *Repository) UpdateLesson(l Lesson) error {
	period := getTstzrange(l.StartTime, l.EndTime)
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE lessons SET subject = $2, summary = $3, tutor = $4, student = $5, scheduled = $6, period = $7 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, l.Id, l.Subject.Id, l.Summary, l.Tutor, l.Student, l.Scheduled, period)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLessonById(id string) (Lesson, error) {
	var period pgtype.Tstzrange
	var l Lesson

	sql := `SELECT lessons.id, subjects.id, subjects.name, subjects.standard, lessons.summary, lessons.tutor, lessons.student, lessons.scheduled, lessons.period FROM lessons INNER JOIN subjects ON lessons.subject = subjects.id WHERE lessons.id = $1`
	if err := r.dbPool.QueryRow(context.Background(), sql, id).Scan(&l.Id, &l.Subject.Id, &l.Subject.Name, &l.Subject.Standard, &l.Summary, &l.Tutor, &l.Student, &l.Scheduled, &period); err != nil {
		return l, err
	}

	period.Upper.AssignTo(&l.EndTime)
	period.Lower.AssignTo(&l.StartTime)
	return l, nil
}
