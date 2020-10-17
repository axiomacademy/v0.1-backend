package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Subject struct {
	Id       string
	Name     string
	Standard string
}

func (r *Repository) ToSubjectModel(sb Subject) model.Subject {
	subject := model.Subject{Name: model.SubjectName(sb.Name), Standard: model.SubjectStandard(sb.Standard)}
	return subject
}

// Creates the subject if it doesn't exist, else retrieves the approprivate subject
func (r *Repository) GetSubject(name string, standard string) (Subject, error) {
	var s Subject

	s.Name = name
	s.Standard = standard

	// Check if it exists first
	sql := `SELECT id FROM subjects WHERE name = $1 AND standard = $2`
	err := r.dbPool.QueryRow(context.Background(), sql, s.Name, s.Standard).Scan(&s.Id)
	if err == pgx.ErrNoRows {
		s.Id = uuid.New()

		tx, err := r.dbPool.Begin(context.Background())
		if err != nil {
			return s, err
		}

		defer tx.Rollback(context.Background())

		sql := `INSERT INTO subjects (id, name, standard) VALUES ($1, $2, $3)`
		_, err = tx.Exec(context.Background(), sql, s.Id, s.Name, s.Standard)

		if err != nil {
			return s, err
		}

		err = tx.Commit(context.Background())
		if err != nil {
			return s, err
		}
	} else if err != nil {
		return s, err
	}

	return s, nil
}

// An utility function to make converting the graphql resovler type easier
func (r *Repository) GetSubjects(subjects []*model.NewSubject) ([]Subject, error) {
	var dbSubjects []Subject
	for _, s := range subjects {
		dbs, err := r.GetSubject(s.Name.String(), s.Standard.String())
		if err != nil {
			return dbSubjects, err
		}

		dbSubjects = append(dbSubjects, dbs)
	}

	return dbSubjects, nil
}

// Get subject by subject UUID
func (r *Repository) GetSubjectById(subid string) (Subject, error) {
	sql := `SELECT id, name, standard FROM subjects WHERE id = $1`

	var subject Subject

	if err := r.dbPool.QueryRow(context.Background(), sql, subid).Scan(&subject.Id, &subject.Name, &subject.Standard); err != nil {
		return subject, err
	}

	return subject, nil
}

// Get all the subjects associated with a tutor based on tutor UUID
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

// Add a fresh set of subjects to the tutor
func (r *Repository) AddSubjectsToTutor(tid string, subids []string) error {
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
