package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

func (r *mutationResolver) CreateStudent(ctx context.Context, input model.NewStudent) (string, error) {
	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	s, err := r.Repo.CreateStudent(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic)
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
	s, err := r.Repo.GetStudentByUsername(input.Username)
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
	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		return "", errors.New("Error hashing password")
	}

	// DEFAULT RATING IS 3
	// Create db.Subject type
	subjects, err := r.Repo.GetSubjects(input.Subjects)
	if err != nil {
		return "", err
	}

	t, err := r.Repo.CreateTutor(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic, input.HourlyRate, 3, input.Bio, input.Education, subjects)
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
	t, err := r.Repo.GetTutorByUsername(input.Username)
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
	t.Status = input.String()
	t.LastSeen = time.Now()

	err = r.Repo.UpdateTutor(t)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		return "", errors.New("Error generating token")
	}

	return token, nil
}

func (r *mutationResolver) MatchOnDemand(ctx context.Context, input model.MatchRequest) (string, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	if utype == "s" {
		s := u.(db.Student)
		subject, err := r.Repo.GetSubject(input.Subject.Name.String(), input.Subject.Standard.String())
		if err != nil {
			return "", err
		}

		mid, err := r.Ms.MatchOnDemand(s, subject)
		return mid, err
	} else if utype == "t" {
		return "", errors.New("Only students can request for matches")
	} else {
		return "", errors.New("Unauthorised, please log in")
	}
}

func (r *mutationResolver) AcceptMatch(ctx context.Context, input string) (*model.Lesson, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		return nil, errors.New("Only students can request for matches")
	} else if utype == "t" {
		t := u.(db.Tutor)
		l, err := r.Ms.AcceptOnDemandMatch(t, input)
		ml, err := r.Repo.ToLessonModel(*l)
		return &ml, err
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}
}

func (r *mutationResolver) UpdateNotification(ctx context.Context, input model.UpdateNotification) (*model.Notification, error) {
	n, err := r.Repo.GetNotificationById(input.ID)
	if err != nil {
		return nil, err
	}

	n.Read = input.Read

	err = r.Repo.UpdateNotification(n)
	if err != nil {
		return nil, err
	}

	mn := r.Repo.ToNotificationModel(n)
	return &mn, nil
}

func (r *mutationResolver) RegisterPushNotification(ctx context.Context, input string) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Self(ctx context.Context) (model.User, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		s := u.(db.Student)
		return r.Repo.ToStudentModel(s), nil
	} else if utype == "t" {
		t := u.(db.Tutor)
		return r.Repo.ToTutorModel(t), nil
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
		dbLessons, err = r.Repo.GetStudentLessons(s.Id)
		if err != nil {
			return nil, err
		}
	} else if utype == "t" {
		t := u.(db.Tutor)
		dbLessons, err = r.Repo.GetTutorLessons(t.Id)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}

	// Convert dbLessons to gql Lesson Type
	var lessons []*model.Lesson
	for _, l := range dbLessons {
		rl, err := r.Repo.ToLessonModel(l)
		if err != nil {
			return nil, err
		}
		lessons = append(lessons, &rl)
	}

	return lessons, nil
}

func (r *queryResolver) Notifications(ctx context.Context, input model.PaginatedRequest) ([]*model.Notification, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var dbNotifications []db.Notification
	if utype == "s" {
		s := u.(db.Student)
		dbNotifications, err = r.Repo.GetUserNotifications(s.Id, input.StartTime, input.EndTime)
		if err != nil {
			return nil, err
		}
	} else if utype == "t" {
		t := u.(db.Tutor)
		dbNotifications, err = r.Repo.GetUserNotifications(t.Id, input.StartTime, input.EndTime)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}

	// Convert dbLessons to gql Lesson Type
	var notifications []*model.Notification
	for _, n := range dbNotifications {
		rn := r.Repo.ToNotificationModel(n)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, &rn)
	}

	return notifications, nil
}

func (r *queryResolver) CheckForMatch(ctx context.Context, input string) (*model.Lesson, error) {
	u, utype, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if utype == "s" {
		s := u.(db.Student)
		l, err := r.Ms.GetOnDemandMatch(s, input)
		ml, err := r.Repo.ToLessonModel(*l)
		return &ml, err
	} else if utype == "t" {
		return nil, errors.New("Only students can request for matches")
	} else {
		return nil, errors.New("Unauthorised, please log in")
	}
}

func (r *subscriptionResolver) SubscribeMatchNotifications(ctx context.Context, user string) (<-chan *model.MatchNotification, error) {
	nchan := r.Ns.CreateUserMatchChannel(user)

	// Delete channel when done
	go func() {
		<-ctx.Done()
		r.Ns.DeleteUserMatchChannel(user)
	}()

	return *nchan, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
