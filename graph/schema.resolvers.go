package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

func (r *mutationResolver) CreateStudent(ctx context.Context, input model.NewStudent) (string, error) {
	s := &db.Student{}

	err := s.Create(input.Email, input.Password, input.ProfilePic)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(s.Id, r.Secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) CreateTutor(ctx context.Context, input model.NewTutor) (string, error) {
	t := &db.Tutor{}

	// DEFAULT RATING IS 3
	err := t.Create(input.Email, input.Password, input.ProfilePic, input.HourlyRate, 3, input.Bio, input.Education, input.Subjects)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (r *mutationResolver) CreateLesson(ctx context.Context, input model.NewLesson) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Self(ctx context.Context) (model.User, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		s := u.(db.Student)
		return s.ToModel(), nil
	} else if utype == "t" {
		t := u.(db.Tutor)
		return t.ToModel(), nil
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}
}

func (r *queryResolver) Lessons(ctx context.Context) ([]*model.Lesson, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var dbLessons []db.Lesson
	if utype == "s" {
		s := u.(db.Student)
		dbLessons, err = s.GetLessons()
		if err != nil {
			return nil, err
		}
	} else if utype == "t" {
		t := u.(db.Tutor)
		dbLessons, err = t.GetLessons()
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}

	// Convert dbLessons to gql Lesson Type
	var lessons []*model.Lesson
	for _, l := range dbLessons {
		rl := l.ToModel()
		lessons = append(lessons, &rl)
	}

	return lessons, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
