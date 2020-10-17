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
	Subjects       []string
	Status         string
	LastSeen       time.Time
	PushToken      string
}

// Convert db.Tutor to model.Tutor
func (r *Repository) ToTutorModel(t Tutor) (model.Tutor, error) {
	var subjects []*model.Subject

	dbSubjects, err := r.getTutorSubjects(t.Id)
	if err != nil {
		return model.Tutor{}, err
	}

	for _, dbSubject := range dbSubjects {
		subject := r.ToSubjectModel(dbSubject)
		subjects = append(subjects, &subject)
	}

	return model.Tutor{ID: t.Id, Username: t.Username, FirstName: t.FirstName, LastName: t.LastName, Email: t.Email, ProfilePic: t.ProfilePic, HourlyRate: t.HourlyRate, Bio: t.Bio, Rating: t.Rating, Education: t.Education, Subjects: subjects}, nil
}

// Creates a new tutor, takes subject IDs
func (r *Repository) CreateTutor(username string, firstName string, lastName string, email string, hashedPassword string, profile_pic string, hourly_rate int, rating int, bio string, education []string, subjects []string) (Tutor, error) {

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
	if err = r.addSubjectsToTutor(t.Id, t.Subjects); err != nil {
		return t, err
	}

	return t, nil
}

// Update tutor to the passed in tutor struct
func (r *Repository) UpdateTutor(t Tutor) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `UPDATE tutors SET first_name = $2, last_name = $3, email = $4, profile_pic = $5, hourly_rate = $6, bio = $7, rating = $8, education = $9, status = $10, last_seen = $11, push_token = $12 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, t.Id, t.FirstName, t.LastName, t.Email, t.ProfilePic, t.HourlyRate, t.Bio, t.Rating, t.Education, t.Status, t.LastSeen, t.PushToken)

	if err != nil {
		return err
	}

	if err = tx.Commit(context.Background()); err != nil {
		return err
	}

	// Updating Subjects to tutor
	if err = r.removeSubjectsFromTutor(t.Id); err != nil {
		return err
	}
	if err = r.addSubjectsToTutor(t.Id, t.Subjects); err != nil {
		return err
	}

	return nil
}

// Get the tutor based on the Tutor UUID
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
	subjects, err := r.getTutorSubjects(t.Id)
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
	subjects, err := r.getTutorSubjects(t.Id)
	if err != nil {
		return t, err
	}

	var subids []string

	for _, subject := range subjects {
		subids = append(subids, subject.Id)
	}

	t.Subjects = subids

	lastSeen.AssignTo(&t.LastSeen)

	return t, nil
}

// Get lessons that the tutor teaches, bounded by a time period
func (r *Repository) GetTutorLessons(tid string, startTime time.Time, endTime time.Time) ([]Lesson, error) {
	sql := `SELECT id, subject, summary, tutor, student, scheduled, period FROM lessons WHERE tutor = $1 and $2 @> period`

	var lessons []Lesson

	period := getTstzrange(startTime, endTime)

	rows, err := r.dbPool.Query(context.Background(), sql, tid, period)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var lesson Lesson
		var period pgtype.Tstzrange
		var sid string

		err := rows.Scan(&lesson.Id, &sid, &lesson.Tutor, &lesson.Student, &lesson.Scheduled, &period)

		if err != nil {
			return nil, err
		}

		period.Upper.AssignTo(&lesson.EndTime)
		period.Lower.AssignTo(&lesson.StartTime)

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

// Get all the subjects associated with a tutor based on tutor UUID
func (r *Repository) getTutorSubjects(tid string) ([]Subject, error) {
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

// Add a fresh set of subjects to the tutor
func (r *Repository) addSubjectsToTutor(tid string, subids []string) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO teaching (tutor, subject) VALUES ($1, $2)`
	for _, s := range subids {
		_, err = tx.Exec(context.Background(), sql, tid, s)
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

// Delete all subjects from the tutor, should do this before you add a fresh set of subjects
func (r *Repository) removeSubjectsFromTutor(tid string) error {
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
