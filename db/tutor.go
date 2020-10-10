package db

import (
	"context"
	"time"

	"github.com/jackc/pgtype"
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
	Subjects       []Subject
	Status         string
	LastSeen       time.Time
	PushToken      string
}

func (r *Repository) ToTutorModel(t Tutor) model.Tutor {
	var subjects []*model.Subject
	for _, dbSubject := range t.Subjects {
		subject := r.ToSubjectModel(dbSubject)
		subjects = append(subjects, &subject)
	}

	return model.Tutor{ID: t.Id, Username: t.Username, FirstName: t.FirstName, LastName: t.LastName, Email: t.Email, ProfilePic: t.ProfilePic, HourlyRate: t.HourlyRate, Bio: t.Bio, Rating: t.Rating, Education: t.Education, Subjects: subjects}
}

func (r *Repository) CreateTutor(username string, firstName string, lastName string, email string, hashedPassword string, profile_pic string, hourly_rate int, rating int, bio string, education []string, subjects []Subject) (Tutor, error) {

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
	t.Subjects = subjects
	t.Status = "UNAVAILABLE"
	t.LastSeen = time.Now()

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return t, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO tutors (id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, rating, bio, education, status, last_seen) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.Username, t.FirstName, t.LastName, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Rating, t.Bio, t.Education, t.Status, t.LastSeen)
	if err != nil {
		return t, err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return t, err
	}

	// Add Subjects to tutor
	if err = r.AddSubjectsToTutor(t.Id, t.Subjects); err != nil {
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

	sql := `UPDATE tutors SET first_name = $2, last_name = $3, email = $4, hashed_password = $5, profile_pic = $6, hourly_rate = $7, bio = $8, rating = $9, education = $10, status = $11, last_seen = $12, push_token = $13 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.FirstName, t.LastName, t.Email, t.HashedPassword, t.ProfilePic, t.HourlyRate, t.Bio, t.Rating, t.Education, t.Status, t.LastSeen, t.PushToken)

	if err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	// Updating Subjects to tutor
	if err = r.RemoveSubjectsFromTutor(t.Id); err != nil {
		return err
	}
	if err = r.AddSubjectsToTutor(t.Id, t.Subjects); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetTutorById(id string) (Tutor, error) {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, status, last_seen, push_token FROM tutors WHERE id = $1`

	var lastSeen pgtype.Timestamptz
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
		&t.Status,
		&lastSeen,
		&t.PushToken); err != nil {
		return t, err
	}

	// Populating subjects separately it needs to be parsed
	subjects, err := r.GetTutorSubjects(t.Id)
	if err != nil {
		return t, err
	}

	lastSeen.AssignTo(&t.LastSeen)
	t.Subjects = subjects

	return t, nil
}

func (r *Repository) GetTutorByUsername(username string) (Tutor, error) {
	sql := `SELECT id, username, first_name, last_name, email, hashed_password, profile_pic, hourly_rate, bio, rating, education, status, last_seen, push_token FROM tutors WHERE username = $1`

	var lastSeen pgtype.Timestamptz
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
		&t.Status,
		&lastSeen,
		&t.PushToken); err != nil {
		return t, err
	}

	// Populating subjects separately because it is an enum array
	subjects, err := r.GetTutorSubjects(t.Id)
	if err != nil {
		return t, err
	}

	lastSeen.AssignTo(&t.LastSeen)
	t.Subjects = subjects
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
		var date pgtype.Timestamptz
		var sid string

		err := rows.Scan(&lesson.Id, &sid, &lesson.Tutor, &lesson.Student, &lesson.Duration, &date, &lesson.Chat)

		if err != nil {
			return nil, err
		}

		date.AssignTo(&lesson.Date)
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

func (r *Repository) GetTutorSubjects(tid string) ([]Subject, error) {
	sql := `SELECT subjects.id, subjects.name, subjects.standard FROM subjects INNER JOIN teaching ON subjects.id = teaching.subject WHERE teaching.tutor = $1`

	var dbSubjects []Subject

	rows, err := r.dbPool.Query(context.Background(), sql, tid)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var subject Subject
		if err := rows.Scan(&subject.Id, &subject.Name, &subject.Standard); err != nil {
			return nil, err
		}

		dbSubjects = append(dbSubjects, subject)

		if err = rows.Err(); err != nil {
			return nil, err
		}
	}

	return dbSubjects, nil
}

func (r *Repository) AddSubjectsToTutor(tid string, subjects []Subject) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO teaching (tutor, subject) VALUES ($1, $2)`
	for _, s := range subjects {
		_, err = tx.Exec(context.Background(), sql, tid, s.Id)
		if err != nil {
			return err
		}
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RemoveSubjectsFromTutor(tid string) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `DELETE FROM teaching WHERE tutor = $1`
	_, err = tx.Exec(context.Background(), sql, tid)
	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}
