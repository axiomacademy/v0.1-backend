package match

import (
	"time"

	"github.com/dgraph-io/badger/v2"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/model"
	"github.com/solderneer/axiom-backend/services/heartbeat"
	"github.com/solderneer/axiom-backend/services/notifs"
)

type MatchService struct {
	db *badger.DB

	Hs   *heartbeat.HeartbeatService
	Ns   *notifs.NotifService
	Repo *db.Repository
}

func (ms *MatchService) Init(badgerDir string) error {
	var err error
	ms.db, err = ms.openBadger(badgerDir)

	return err
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

func (ms *MatchService) MatchOnDemand(s db.Student, subject string, subject_level string) error {
	// Need to first collect top 10 students, ordered by affinity, send match notifications to each of them (timeout 20 seconds), once a match is found set match Id to a created lesson

	go func() {
		tids := []string{"1", "2", "3", "4", "5"}

		for _, tid := range tids {
			n := model.Notification{
				student:       ms.Repo.ToStudentModel(s),
				subject:       subject,
				subject_level: subject_level,
				expiry:        time.Now().Add(time.Duration(30) * time.Second),
			}

			ms.Ns.SendNotification(n, tid)
			time.Sleep(time.Duration(30) * time.Second)

			// Go and check the match queue for a match
		}
	}()

	return nil
}

func (ms *MatchService) AcceptOnDemandMatch(t db.Tutor) error {
	return nil
}

func (ms *MatchService) GetOnDemandMatch(s db.Student) error {
	return nil
}
