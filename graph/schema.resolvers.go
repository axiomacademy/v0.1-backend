package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/utilities/auth"
)

func (r *mutationResolver) CreateStudent(ctx context.Context, input model.NewStudent) (string, error) {
	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		r.sendError(err, "Cannot hash password")
		return "", InternalServerError
	}

	s, err := r.Repo.CreateStudent(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic)
	if err != nil {
		return "", err
	}

	token, err := auth.GenerateToken(s.Id, r.Secret)
	if err != nil {
		r.sendError(err, "Cannot generate JWT")
		return "", InternalServerError
	}

	return token, nil
}

func (r *mutationResolver) LoginStudent(ctx context.Context, input model.LoginInfo) (string, error) {
	s, err := r.Repo.GetStudentByUsername(input.Username)
	if err != nil {
		r.Logger.WithFields(log.Fields{
			"err": err.Error(),
		}).Error("Cannot log in")
		return "", errors.New("Invalid username")
	}

	ok := auth.CheckPasswordHash(input.Password, s.HashedPassword)
	if !ok {
		return "", errors.New("Invalid password")
	}

	token, err := auth.GenerateToken(s.Id, r.Secret)
	if err != nil {
		r.sendError(err, "Cannot generate JWT")
		return "", InternalServerError
	}

	return token, nil
}

func (r *mutationResolver) CreateTutor(ctx context.Context, input model.NewTutor) (string, error) {
	// Hashing password
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		r.sendError(err, "Cannot hash password")
		return "", InternalServerError
	}

	// DEFAULT RATING IS 3
	// Create db.Subject type
	subjects, err := r.Repo.GetSubjects(input.Subjects)
	if err != nil {
		r.sendError(err, "Cannnot retrieve subjects from database")
		return "", InternalServerError
	}

	var subids []string
	for _, subject := range subjects {
		subids = append(subids, subject.Id)
	}

	t, err := r.Repo.CreateTutor(input.Username, input.FirstName, input.LastName, input.Email, hashedPassword, input.ProfilePic, input.HourlyRate, 3, input.Bio, input.Education, subids)
	if err != nil {
		r.sendError(err, "Cannot create tutor in database")
		return "", InternalServerError
	}

	token, err := auth.GenerateToken(t.Id, r.Secret)
	if err != nil {
		r.sendError(err, "Cannot generate JWT")
		return "", InternalServerError
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
		r.sendError(err, "Cannot generate JWT")
		return "", InternalServerError
	}

	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context) (string, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", Unauthorised
	}

	var token string

	switch user := u.(type) {
	case db.Student:
		token, err = auth.GenerateToken(user.Id, r.Secret)
		if err != nil {
			r.sendError(err, "Cannot generate JWT")
			return "", InternalServerError
		}
	case db.Tutor:
		token, err = auth.GenerateToken(user.Id, r.Secret)
		if err != nil {
			r.sendError(err, "Cannot generate JWT")
			return "", InternalServerError
		}
	default:
		return "", Unauthorised
	}

	return token, nil
}

func (r *mutationResolver) UpdateHeartbeat(ctx context.Context, input model.HeartbeatStatus) (string, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	switch user := u.(type) {
	case db.Student:
		r.sendError(err, "Only tutors have persission to update heartbeat")
		return "", Unauthorised
	case db.Tutor:
		user.Status = input.String()
		user.LastSeen = time.Now()

		err = r.Repo.UpdateTutor(user)
		if err != nil {
			return "", err
		}

		token, err := auth.GenerateToken(user.Id, r.Secret)
		if err != nil {
			r.sendError(err, "Cannot generate JWT")
			return "", InternalServerError
		}

		return token, nil
	default:
		return "", Unauthorised
	}
}

func (r *mutationResolver) RequestOnDemandMatch(ctx context.Context, input model.OnDemandMatchRequest) (string, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	switch user := u.(type) {
	case db.Student:
		subject, err := r.Repo.GetSubject(input.Subject.Name.String(), input.Subject.Standard.String())
		if err != nil {
			r.sendError(err, "Cannot get subject from database")
			return "", InternalServerError
		}

		mid, err := r.Ms.MatchOnDemand(user, subject, 20)
		return mid, err
	case db.Tutor:
		r.sendError(err, "Only students can request for a match")
		return "", Unauthorised
	default:
		return "", Unauthorised
	}
}

func (r *mutationResolver) RequestScheduledMatch(ctx context.Context, input model.ScheduledMatchRequest) (string, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	switch user := u.(type) {
	case db.Student:
		// Retrieve the subject
		sub, err := r.Repo.GetSubject(input.Subject.Name.String(), input.Subject.Standard.String())
		if err != nil {
			r.sendError(err, "Cannot retrieve subject from db")
			return "", InternalServerError
		}
		// Retrieve the tutor
		t, err := r.Repo.GetTutorById(input.Tutor)
		if err != nil {
			r.sendError(err, "Cannot retrieve tutor from db")
			return "", InternalServerError
		}
		m, err := r.Ms.RequestScheduledMatch(user, t, sub, input.Time.StartTime, input.Time.EndTime)
		if err != nil {
			r.sendError(err, "Cannot request new match")
			return "", InternalServerError
		}

		return m.Id, nil
	case db.Tutor:
		r.sendError(err, "Only students can request for a match")
		return "", Unauthorised
	default:
		return "", Unauthorised
	}
}

func (r *mutationResolver) AcceptOnDemandMatch(ctx context.Context, input string) (*model.Lesson, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		r.sendError(err, "Only tutors can accept a match")
		return nil, Unauthorised
	case db.Tutor:
		l, err := r.Ms.AcceptOnDemandMatch(input, user)
		ml, err := r.Repo.ToLessonModel(l)
		if err != nil {
			r.sendError(err, "Cannot accept match")
			return nil, InternalServerError
		}
		return &ml, nil
	default:
		return nil, Unauthorised
	}
}

func (r *mutationResolver) AcceptScheduledMatch(ctx context.Context, input string) (*model.Lesson, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		r.sendError(err, "Only tutors can accept a match")
		return nil, Unauthorised
	case db.Tutor:
		l, err := r.Ms.AcceptScheduledMatch(input, user)
		ml, err := r.Repo.ToLessonModel(l)
		if err != nil {
			r.sendError(err, "Cannot accept match")
			return nil, InternalServerError
		}
		return &ml, nil
	default:
		return nil, Unauthorised
	}
}

func (r *mutationResolver) UpdateNotification(ctx context.Context, input model.UpdateNotification) (*model.Notification, error) {
	n, err := r.Repo.GetNotificationById(input.ID)
	if err != nil {
		r.sendError(err, "Cannot retrieve notification from database")
		return nil, InternalServerError
	}

	// Check that notification is for the correct user
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		if n.Student != user.Id {
			return nil, Unauthorised
		}
	case db.Tutor:
		if n.Tutor != user.Id {
			return nil, Unauthorised
		}
	default:
		return nil, Unauthorised
	}

	n.Read = input.Read

	err = r.Repo.UpdateNotification(n)
	if err != nil {
		r.sendError(err, "Cannot update notification in database")
		return nil, InternalServerError
	}

	mn := r.Repo.ToNotificationModel(n)
	return &mn, nil
}

func (r *mutationResolver) RegisterPushNotification(ctx context.Context, input string) (string, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return "", err
	}

	switch user := u.(type) {
	case db.Student:
		user.PushToken = input
		err := r.Repo.UpdateStudent(user)

		if err != nil {
			r.sendError(err, "Cannot update student in db")
			return "", InternalServerError
		}
		return input, nil
	case db.Tutor:
		user.PushToken = input
		err := r.Repo.UpdateTutor(user)

		if err != nil {
			r.sendError(err, "Cannot update tutor in db")
			return "", InternalServerError
		}
		return input, nil
	default:
		return "", Unauthorised
	}
}

func (r *queryResolver) Self(ctx context.Context) (model.User, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		return r.Repo.ToStudentModel(user), nil
	case db.Tutor:
		return r.Repo.ToTutorModel(user)
	default:
		return nil, Unauthorised
	}
}

func (r *queryResolver) Lessons(ctx context.Context, input model.TimeRangeRequest) ([]*model.Lesson, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var dbLessons []db.Lesson

	switch user := u.(type) {
	case db.Student:
		dbLessons, err = r.Repo.GetStudentLessons(user.Id, input.StartTime, input.EndTime)
		if err != nil {
			r.sendError(err, "Cannot retrieve student lessons from db")
			return nil, InternalServerError
		}
	case db.Tutor:
		dbLessons, err = r.Repo.GetTutorLessons(user.Id, input.StartTime, input.EndTime)
		if err != nil {
			r.sendError(err, "Cannot retrieve tutor lessons from db")
			return nil, InternalServerError
		}
	default:
		return nil, Unauthorised
	}

	// Convert dbLessons to gql Lesson Type
	var lessons []*model.Lesson
	for _, l := range dbLessons {
		rl, err := r.Repo.ToLessonModel(l)
		if err != nil {
			r.sendError(err, "Cannot parse lesson")
			return nil, InternalServerError
		}
		lessons = append(lessons, &rl)
	}

	return lessons, nil
}

func (r *queryResolver) PendingMatches(ctx context.Context) ([]*model.Match, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var dbMatches []db.Match

	switch user := u.(type) {
	case db.Student:
		dbMatches, err = r.Repo.GetStudentPendingMatches(user.Id)
		if err != nil {
			r.sendError(err, "Cannot retrieve student pending matches from database")
			return nil, InternalServerError
		}
	case db.Tutor:
		dbMatches, err = r.Repo.GetTutorPendingMatches(user.Id)
		if err != nil {
			r.sendError(err, "Cannot retrieve tutor pending matches from database")
			return nil, InternalServerError
		}
	default:
		return nil, Unauthorised
	}

	var modelMatches []*model.Match
	for _, match := range dbMatches {
		model, err := r.Repo.ToMatchModel(match)
		if err != nil {
			r.sendError(err, "Cannot cast dbmodel to gql model")
			return nil, InternalServerError
		}
		modelMatches = append(modelMatches, &model)
	}

	return modelMatches, nil
}

func (r *queryResolver) Notifications(ctx context.Context, input model.TimeRangeRequest) ([]*model.Notification, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var dbNotifications []db.Notification

	switch user := u.(type) {
	case db.Student:
		dbNotifications, err = r.Repo.GetUserNotifications(user.Id, input.StartTime, input.EndTime)
		if err != nil {
			r.sendError(err, "Cannot retrieve student notifications from database")
			return nil, InternalServerError
		}
	case db.Tutor:
		dbNotifications, err = r.Repo.GetUserNotifications(user.Id, input.StartTime, input.EndTime)
		if err != nil {
			r.sendError(err, "Cannot retrieve tutor notifications from database")
			return nil, InternalServerError
		}
	default:
		return nil, Unauthorised
	}

	// Convert dbNotifications to gql Notifications Type
	var notifications []*model.Notification
	for _, n := range dbNotifications {
		rn := r.Repo.ToNotificationModel(n)
		if err != nil {
			r.sendError(err, "Cannot parse notification model")
			return nil, InternalServerError
		}
		notifications = append(notifications, &rn)
	}

	return notifications, nil
}

func (r *queryResolver) GetScheduledMatches(ctx context.Context, input model.ScheduledMatchParameters) ([]*model.Tutor, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		subject, err := r.Repo.GetSubject(input.Subject.Name.String(), input.Subject.Standard.String())
		if err != nil {
			r.sendError(err, "Cannot get subject from database")
			return nil, InternalServerError
		}

		tids, err := r.Ms.MatchScheduled(user, subject, input.Time.StartTime, input.Time.EndTime, 20)
		if err != nil {
			return nil, InternalServerError
		}

		var tutors []*model.Tutor
		for _, tid := range tids {
			dbTutor, err := r.Repo.GetTutorById(tid)
			if err != nil {
				r.sendError(err, "Cannot get tutor from database")
				return nil, InternalServerError
			}

			tutor, err := r.Repo.ToTutorModel(dbTutor)
			if err != nil {
				r.sendError(err, "Cannot parse tutor from database")
				return nil, InternalServerError
			}

			tutors = append(tutors, &tutor)
		}

		return tutors, err
	case db.Tutor:
		r.sendError(err, "Only students can request for a match")
		return nil, Unauthorised
	default:
		return nil, Unauthorised
	}
}

func (r *queryResolver) CheckForMatch(ctx context.Context, input string) (*model.Lesson, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		l, err := r.Ms.CheckForMatch(user, input)
		if err == errors.New("No match found") {
			return nil, err
		} else if err != nil {
			r.sendError(err, "Error retrieving matct")
			return nil, InternalServerError
		}

		ml, err := r.Repo.ToLessonModel(*l)

		if err != nil {
			return &ml, err
		}

		return &ml, nil
	case db.Tutor:
		r.sendError(err, "Only students can long poll for a match")
		return nil, Unauthorised
	default:
		return nil, Unauthorised
	}
}

func (r *subscriptionResolver) SubscribeMatchNotifications(ctx context.Context) (<-chan *model.MatchNotification, error) {
	u, err := auth.UserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	switch user := u.(type) {
	case db.Student:
		r.sendError(err, "Only tutors can subscribe to match notifications")
		return nil, Unauthorised
	case db.Tutor:
		nchan := r.Ns.CreateUserMatchChannel(user.Id)

		// Delete channel when done
		go func() {
			<-ctx.Done()
			r.Ns.DeleteUserMatchChannel(user.Id)
		}()

		return *nchan, nil
	default:
		return nil, Unauthorised
	}
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
