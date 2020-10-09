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

func (r *Repository) GetSubjectById(sid string) (Subject, error) {
	sql := `SELECT id, name, standard FROM subjects WHERE id = $1`

	var subject Subject

	if err := r.dbPool.QueryRow(context.Background(), sql, sid).Scan(&subject.Id, &subject.Name, &subject.Standard); err != nil {
		return subject, err
	}

	return subject, nil
}
