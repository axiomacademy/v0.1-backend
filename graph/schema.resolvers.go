package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/heartbeat"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

func (r *mutationResolver) CreateStudent(ctx context.Context, input model.NewStudent) (string, error) {
	s := &db.Student{}

	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	err = s.Create(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(s.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
	}

	return token, nil
}

func (r *mutationResolver) LoginStudent(ctx context.Context, input model.LoginInfo) (string, error) {
	s := &db.Student{}

	err := s.GetByUsername(input.Username)
	if err != nil {
		return "", errors.New("Invalid username")
	}

	ok := auth.CheckPasswordHash(input.Password, s.HashedPassword)
	if !ok {
		return "", errors.New("Invalid password")
	}

	token, err := auth.GenerateToken(s.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
	}

	return token, nil
}

func (r *mutationResolver) CreateTutor(ctx context.Context, input model.NewTutor) (string, error) {
	t := &db.Tutor{}

	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	// DEFAULT RATING IS 3
	err = t.Create(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic, input.HourlyRate, 3, input.Bio, input.Education, input.Subjects)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
	}

	return token, nil
}

func (r *mutationResolver) LoginTutor(ctx context.Context, input model.LoginInfo) (string, error) {
	t := &db.Tutor{}

	err := t.GetByUsername(input.Username)
	if err != nil {
		return "", errors.New("Invalid username")
	}

	ok := auth.CheckPasswordHash(input.Password, t.HashedPassword)
	if !ok {
		return "", errors.New("Invalid password")
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
	}

	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	var token string

	if utype == "s" {
		s := u.(db.Student)

		token, err = auth.GenerateToken(s.Id, r.Secret)
		if err != nil {
			return "", errors.New("Error generating token")
		}
	} else if utype == "t" {
		t := u.(db.Tutor)

		token, err = auth.GenerateToken(t.Id, r.Secret)
		if err != nil {
			return "", errors.New("Error generating token")
		}
	} else {
		return "", errors.New("Unauthorised, please log in")
	}

	return token, nil
}

func (r *mutationResolver) UpdateHeartbeat(ctx context.Context, input model.HeartbeatStatus) (string, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	if utype != "t" {
		return "", errors.New("Invalid user type for Heartbeat")
	}

	t := u.(db.Tutor)

	err = heartbeat.SetHeartbeat(t.Id, input)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
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
