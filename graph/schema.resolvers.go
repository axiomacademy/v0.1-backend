package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
)

func (r *mutationResolver) CreateStudent(ctx context.Context, input model.NewStudent) (*model.Student, error) {
	student := &model.Student{
		ID:             input.ID,
		Email:          input.Email,
		HashedPassword: input.HashedPassword,
		Lessons:        []*model.Lesson{},
	}

	r.students = append(r.students, student)
	return student, nil
}

func (r *queryResolver) Self(ctx context.Context) (model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Users(ctx context.Context) ([]model.User, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Lessons(ctx context.Context) ([]*model.Lesson, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
