package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

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

func (r *queryResolver) Self(ctx context.Context) (model.User, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		s := u.(db.Student)
		return &model.Student{ID: s.Id, Email: s.Email, ProfilePic: s.ProfilePic}, nil
	} else if utype == "t" {
		t := u.(db.Tutor)
		return &model.Tutor{ID: t.Id, Email: t.Email, ProfilePic: t.ProfilePic, HourlyRate: t.HourlyRate, Bio: t.Bio, Rating: t.Rating, Education: t.Education, Subjects: t.Subjects}, nil
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}
}

func (r *queryResolver) Lessons(ctx context.Context) ([]*model.Lesson, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		s := u.(db.Student)
		dbLessons, err := s.GetLessons()
		if err != nil {
			return nil, err
		}

		// Convert dbLessons to gql Lesson Type
		var lessons []*model.Lesson
		for _, l := range dbLessons {
			lessons = append(lessons, &model.Lesson{ID: l.Id, Subject: l.Subject, Summary: l.Summary, Tutor: l.Tutor, Student: l.Student, Duration: l.Duration, Date: l.Date, Chat: l.Chat})
		}

		return lessons, nil
	} else if utype == "t" {
		t := u.(db.Tutor)
		dbLessons, err := t.GetLessons()
		if err != nil {
			return nil, err
		}

		// Convert dbLessons to gql Lesson Type
		var lessons []*model.Lesson
		for _, l := range dbLessons {
			lessons = append(lessons, &model.Lesson{ID: l.Id, Subject: l.Subject, Summary: l.Summary, Tutor: l.Tutor, Student: l.Student, Duration: l.Duration, Date: l.Date, Chat: l.Chat})
		}

		return lessons, nil
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
