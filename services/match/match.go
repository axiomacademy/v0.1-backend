package match

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	log "github.com/sirupsen/logrus"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/services/notifs"
)

type MatchService struct {
	logger *log.Logger

	secret string
	ns     *notifs.NotifService
	repo   *db.Repository
}

func (ms *MatchService) Init(logger *log.Logger, secret string, ns *notifs.NotifService, repo *db.Repository) {
	ms.logger = logger
	ms.secret = secret
	ms.ns = ns
	ms.repo = repo

	ms.logger.WithField("service", "match").Info("Successfully initialised")
}

// Retrieves top scheduled matchs based on availability
// Takes in a student, subject, start and end times, plus a limit integer of how many matches to return
func (ms *MatchService) MatchScheduled(s db.Student, subject db.Subject, startTime time.Time, endTime time.Time, limit int) ([]string, error) {
	affinitytids, err := ms.repo.GetAvailableTutors(s.Id, subject.Id, startTime, endTime)
	if err != nil {
		ms.sendError(err, "Cannot retrieve affinity ordered available tutors")
		return nil, err
	}

	var availabletids []string
	var count int
	// Filter affinity Ids
	for _, tid := range affinitytids {

		if count >= limit {
			break
		}

		available, err := ms.repo.CheckTutorAvailability(tid, startTime, endTime)
		if err != nil {
			ms.sendError(err, "Cannot check tutor availability")
			continue
		}

		if available {
			availabletids = append(availabletids, tid)
			count++
		}
	}

	// Fetch more random tutors if there aren't enough affinity matches. Current BATCH_SIZE = 100
	for len(availabletids) < limit {
		randomtids, err := ms.repo.GetRandomAvailableTutors(subject.Id, startTime, endTime, 100)
		if err != nil {
			ms.sendError(err, "Cannot retrieve random available tutors")
			return nil, err
		}

		for _, tid := range randomtids {

			if count >= 20 {
				break
			}

			available, err := ms.repo.CheckTutorAvailability(tid, startTime, endTime)
			if err != nil {
				ms.sendError(err, "Cannot check tutor availability")
				continue
			}

			if available {
				availabletids = append(availabletids, tid)
				count++
			}
		}
	}

	return availabletids, nil
}

func (ms *MatchService) RequestScheduledMatch(s db.Student, t db.Tutor, subject db.Subject, startTime time.Time, endTime time.Time) (db.Match, error) {
	// Create the match
	m, err := ms.repo.CreateScheduledMatch("MATCHING", s.Id, t.Id, subject.Id, startTime, endTime)
	if err != nil {
		ms.sendError(err, "Cannot create match in database")
		return m, err
	}

	// Send a push notification
	n, err := ms.repo.CreateNotification(t.Id, "New scheduled lesson request!", "You have received a new match request from "+s.FirstName, "")
	if err != nil {
		ms.sendError(err, "Cannot create notification in database")
		return m, err
	}

	err = ms.ns.SendPushNotification(n, t.PushToken)
	if err != nil {
		ms.sendError(err, "Cannot send firebase push notification")
		return m, err
	}

	// Handles expiring after one day
	go func() {
		time.Sleep(time.Hour * 24)

		// Send failure push nnotification
		n, err := ms.repo.CreateNotification(s.Id, "Match failed", "Your scheduled match with "+t.FirstName+" has expired", "")
		if err != nil {
			ms.sendError(err, "Cannot create notification in database")
			return
		}

		err = ms.ns.SendPushNotification(n, s.PushToken)
		if err != nil {
			ms.sendError(err, "Cannot send firebase push notification")
			return
		}

		// Update match status
		m.Status = "FAILED"
		err = ms.repo.UpdateMatch(m)
		if err != nil {
			ms.sendError(err, "Cannot update match in database")
			return
		}
	}()

	return m, nil
}

func (ms *MatchService) AcceptScheduledMatch(mid string, t db.Tutor) (db.Lesson, error) {
	var l db.Lesson

	m, err := ms.repo.GetMatchById(mid)
	if err != nil {
		ms.sendError(err, "Cannot retrieve match from database")
		return l, err
	}

	// Check that the tutor in the match and the request tutor match
	if m.Tutor != t.Id {
		ms.sendError(err, "Invalid auth access!")
		return l, errors.New("Unauthorised to access match")
	}

	// Fetch the subject
	sub, err := ms.repo.GetSubjectById(m.Subject)
	if err != nil {
		ms.sendError(err, "Cannot retrieve subject from database")
		return l, err
	}

	// Create the lesson
	l, err = ms.repo.CreateLesson(sub, m.Tutor, m.Student, true, m.StartTime, m.EndTime)
	if err != nil {
		ms.sendError(err, "Cannot create lesson in database")
		return l, err
	}

	// Update match status
	m.Status = "MATCHED"
	m.Lesson = l.Id
	err = ms.repo.UpdateMatch(m)
	if err != nil {
		ms.sendError(err, "Cannot update match in database")
		return l, err
	}

	// Send push notification to student
	s, err := ms.repo.GetStudentById(m.Student)
	if err != nil {
		ms.sendError(err, "Cannot retrieve student from database")
		return l, err
	}

	n, err := ms.repo.CreateNotification(s.Id, "Scheduled lesson confirmed!", "Successfully matched you with "+t.FirstName, "")
	if err != nil {
		ms.sendError(err, "Cannot create notification in database")
		return l, err
	}

	err = ms.ns.SendPushNotification(n, s.PushToken)
	if err != nil {
		ms.sendError(err, "Error sending firebase push notification")
		return l, err
	}

	return l, err
}

// Collects top students, ordered by affinity, send match notifications to each of them (timeout 20 seconds), once a match is found set match Id to a created lesson
// Retrieves all the top  matches for on demand. Limit integer defines how many matches to generate
func (ms *MatchService) MatchOnDemand(s db.Student, subject db.Subject, limit int) (string, error) {
	// Generate match id
	m, err := ms.repo.CreateOnDemandMatch("MATCHING", s.Id, subject.Id)

	if err != nil {
		ms.sendError(err, "Cannot create match in database")
		return "", err
	}

	go func() {
		tids, err := ms.repo.GetOnlineAffinityMatches(s.Id, subject)
		if err != nil {
			m.Status = "FAILED"
			ms.sendError(err, "Error retrieving database matches")

			err := ms.repo.UpdateMatch(m)
			ms.sendError(err, "Updating match to failed")
			return
		}

		if len(tids) < limit {
			rtids, err := ms.repo.GetOnlineRandomMatches(subject, limit-len(tids))
			if err != nil {
				m.Status = "FAILED"
				ms.sendError(err, "Error retrieving database matches")

				err := ms.repo.UpdateMatch(m)
				ms.sendError(err, "Updating match to failed")
				return
			}
			tids = append(tids, rtids...)
		}

		mstudent := ms.repo.ToStudentModel(s)
		token, err := ms.generateMatchToken(m.Id)

		// Error generating token
		if err != nil {
			m.Status = "FAILED"
			ms.sendError(err, "Error generating match token")

			err := ms.repo.UpdateMatch(m)
			ms.sendError(err, "Updating match to failed")
			return
		}

		msubject := ms.repo.ToSubjectModel(subject)

		for _, tid := range tids {
			n := model.MatchNotification{
				Student: &mstudent,
				Subject: &msubject,
				Token:   token,
			}

			ms.ns.SendMatchNotification(n, tid)
			time.Sleep(time.Duration(30) * time.Second)

			// Go and check the match queue for a match
			latestMatch, err := ms.repo.GetMatchById(m.Id)
			if err != nil {
				m.Status = "FAILED"
				ms.sendError(err, "Error retrieving database match")

				err := ms.repo.UpdateMatch(m)
				ms.sendError(err, "Updating match to failed")
				return
			}

			if latestMatch.Status == "MATCHED" {
				// Stop looping
				return
			}
		}

		// No matches found at the moment
		m.Status = "FAILED"
		err = ms.repo.UpdateMatch(m)
		if err != nil {
			ms.sendError(err, "Cannot update database match")
			return
		}

		return
	}()

	return m.Id, err
}

func (ms *MatchService) AcceptOnDemandMatch(token string, t db.Tutor) (db.Lesson, error) {
	var l db.Lesson
	mid, err := ms.parseMatchToken(token)
	if err != nil {
		ms.sendError(err, "Unable to parse match token")
		return l, err
	}

	// Fetching the match
	m, err := ms.repo.GetMatchById(mid)
	if err != nil {
		ms.sendError(err, "Unable to retrieve match")
		return l, err
	}

	// Fetch the match subject
	sub, err := ms.repo.GetSubjectById(m.Subject)
	if err != nil {
		ms.sendError(err, "Unable to retrieve/create subject")
		return l, err
	}

	l, err = ms.repo.CreateLesson(sub, t.Id, m.Student, false, time.Now(), time.Now().Add(15*time.Minute))
	if err != nil {
		ms.sendError(err, "Unable to create lesson in database")
		return l, err
	}

	// Updating match queue
	m.Lesson = l.Id
	m.Status = "MATCHED"
	err = ms.repo.UpdateMatch(m)
	if err != nil {
		ms.sendError(err, "Unable to update match status")
		return l, err
	}

	return l, nil
}

func (ms *MatchService) GetOnDemandMatch(s db.Student, mid string) (*db.Lesson, error) {

	// Fetching the match
	m, err := ms.repo.GetMatchById(mid)
	if err != nil {
		ms.sendError(err, "Unable to retrieve match")
		return nil, err
	}

	// Check that the student is authorised
	if m.Student != s.Id {
		ms.sendError(err, "Invalid auth access!")
		return nil, errors.New("Unauthorised to access match")
	}

	switch m.Status {
	case "MATCHING":
		return nil, errors.New("Still matching")
	case "FAILED":
		return nil, errors.New("Matching failed")
	case "MATCHED":
		l, err := ms.repo.GetLessonById(m.Lesson)
		if err != nil {
			ms.sendError(err, "Unable to retrieve lesson from database")
			return nil, err
		}
		return &l, nil
	}

	return nil, err
}

// GenerateToken generates a jwt token and assign a ,id to it's claims and return it
func (ms *MatchService) generateMatchToken(id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	/* Create a map to store our claims */
	claims := token.Claims.(jwt.MapClaims)
	/* Set token claims */
	claims["id"] = id
	claims["exp"] = time.Now().Add(time.Second * 30).Unix()
	tokenString, err := token.SignedString([]byte(ms.secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

//ParseToken parses a jwt token and returns the mid it it's claims
func (ms *MatchService) parseMatchToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(ms.secret), nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		id := claims["id"].(string)
		return id, nil
	} else {
		return "", err
	}
}

// Making sending errors easier
func (ms *MatchService) sendError(err error, message string) {
	ms.logger.WithFields(log.Fields{
		"service": "match",
		"err":     err.Error(),
	}).Error(message)
}
