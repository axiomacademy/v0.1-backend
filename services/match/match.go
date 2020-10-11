package match

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dgraph-io/badger/v2"
	"github.com/dgrijalva/jwt-go"
	"github.com/pborman/uuid"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/services/notifs"
)

type MatchStatus struct {
	Status          string `json:"status"`
	Sid             string `json:"student_id"`
	SubjectName     string `json:"subject_name"`
	SubjectStandard string `json:"subject_standard"`
	Lid             string `json:"lesson_id"`
}

type MatchService struct {
	logger *log.Logger
	db     *badger.DB

	secret string
	ns     *notifs.NotifService
	repo   *db.Repository
}

func (ms *MatchService) Init(logger *log.Logger, badgerDir string, secret string, ns *notifs.NotifService, repo *db.Repository) {
	ms.logger = logger
	ms.secret = secret
	ms.ns = ns
	ms.repo = repo

	var err error
	ms.db, err = ms.openBadger(badgerDir)
	if err != nil {
		ms.logger.WithField("service", "match").Fatal("Unable to open badger store")
	}

	ms.logger.WithField("service", "match").Info("Successfully initialised")
}

func (ms *MatchService) Close() {
	ms.db.Close()
}

func (ms *MatchService) openBadger(badgerDir string) (*badger.DB, error) {
	if badgerDir == "" {
		return badger.Open(badger.DefaultOptions("").WithInMemory(true))
	} else {
		return badger.Open(badger.DefaultOptions(badgerDir))
	}
}

func (ms *MatchService) updateMatch(mid string, status MatchStatus) error {
	err := ms.db.Update(func(txn *badger.Txn) error {
		raw, err := json.Marshal(status)
		if err != nil {
			ms.sendError(err, "Error marshalling match status")
			return err
		}

		err = txn.Set([]byte(mid), raw)
		if err != nil {
			ms.sendError(err, "Error updating match status")
			return err
		}

		return nil
	})

	return err
}

func (ms *MatchService) deleteMatch(mid string) error {
	err := ms.db.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(mid))
		if err != nil {
			ms.sendError(err, "Error deleting match status")
			return err
		}

		return nil
	})

	return err
}

func (ms *MatchService) getMatch(mid string) (MatchStatus, error) {
	var status MatchStatus

	err := ms.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(mid))
		if err != nil {
			ms.sendError(err, "Error retrieving match status")
			return err
		}

		rawStatus, err := item.ValueCopy(nil)
		if err != nil {
			ms.sendError(err, "Error parsing match status")
			return err
		}

		err = json.Unmarshal([]byte(rawStatus), &status)
		if err != nil {
			ms.sendError(err, "Error marshalling match status")
			return err
		}
		return nil
	})

	return status, err
}

// Need to first collect top 10 students, ordered by affinity, send match notifications to each of them (timeout 20 seconds), once a match is found set match Id to a created lesson
func (ms *MatchService) MatchOnDemand(s db.Student, subject db.Subject) (string, error) {
	// Generate match id
	mid := uuid.New()
	mstatus := MatchStatus{
		Status:          "MATCHING",
		Sid:             s.Id,
		SubjectName:     subject.Name,
		SubjectStandard: subject.Standard,
		Lid:             "",
	}

	// Storing the match on the match queue
	err := ms.updateMatch(mid, mstatus)

	go func() {
		tids, err := ms.repo.GetOnlineAffinityMatches(s.Id, subject)
		if err != nil {
			mstatus := MatchStatus{
				Status: "FAILED",
				Sid:    "",
				Lid:    "",
			}

			ms.sendError(err, "Error retrieving database matches")
			ms.updateMatch(mid, mstatus)
			return
		}

		if len(tids) < 15 {
			rtids, err := ms.repo.GetOnlineRandomMatches(subject, 15-len(tids))
			if err != nil {
				mstatus := MatchStatus{
					Status: "FAILED",
					Sid:    "",
					Lid:    "",
				}

				ms.sendError(err, "Error retrieving database matches")
				ms.updateMatch(mid, mstatus)
				return
			}
			tids = append(tids, rtids...)
		}

		mstudent := ms.repo.ToStudentModel(s)
		token, err := ms.generateMatchToken(mid)

		// Error generating token
		if err != nil {
			mstatus := MatchStatus{
				Status: "FAILED",
				Sid:    "",
				Lid:    "",
			}

			ms.sendError(err, "Error generating match token")
			ms.updateMatch(mid, mstatus)
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
			match, err := ms.getMatch(mid)
			if err != nil {
				mstatus := MatchStatus{
					Status: "FAILED",
					Sid:    "",
					Lid:    "",
				}

				ms.sendError(err, "Error retreiving match")
				ms.updateMatch(mid, mstatus)
				return
			}
			if match.Status == "MATCHED" {
				// Stop looping
				return
			}
		}
		mstatus := MatchStatus{
			Status: "FAILED",
			Sid:    "",
			Lid:    "",
		}

		ms.updateMatch(mid, mstatus)
		return
	}()

	return mid, err
}

func (ms *MatchService) AcceptOnDemandMatch(t db.Tutor, token string) (*db.Lesson, error) {
	mid, err := ms.parseMatchToken(token)
	if err != nil {
		ms.sendError(err, "Unable to parse match token")
		return nil, err
	}

	// Fetching the match
	status, err := ms.getMatch(mid)
	if err != nil {
		ms.sendError(err, "Unable to retrieve match")
		return nil, err
	}

	// Creating the lesson
	subject, err := ms.repo.GetSubject(status.SubjectName, status.SubjectStandard)
	if err != nil {
		ms.sendError(err, "Unable to create subject in database")
		return nil, err
	}

	l, err := ms.repo.CreateLesson(subject, t.Id, status.Sid, 0, time.Now())
	if err != nil {
		ms.sendError(err, "Unable to create lesson in database")
		return nil, err
	}

	// Updating match queue
	status.Lid = l.Id
	status.Status = "MATCHED"
	err = ms.updateMatch(mid, status)
	if err != nil {
		ms.sendError(err, "Unable to update match status")
		return nil, err
	}

	return &l, nil
}

func (ms *MatchService) GetOnDemandMatch(s db.Student, mid string) (*db.Lesson, error) {

	// Fetching the match
	status, err := ms.getMatch(mid)
	if err != nil {
		ms.sendError(err, "Unable to retrieve match")
		return nil, err
	}

	// Check that the student is authorised
	if status.Sid != s.Id {
		ms.sendError(err, "Invalid auth access!")
		return nil, errors.New("Unauthorised to access match")
	}

	switch status.Status {
	case "MATCHING":
		return nil, errors.New("Still matching")
	case "FAILED":
		return nil, errors.New("Matching failed")
	case "MATCHED":
		l, err := ms.repo.GetLessonById(status.Lid)
		if err != nil {
			ms.sendError(err, "Unable to retrieve lesson from database")
			return nil, err
		}
		ms.deleteMatch(mid)
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
