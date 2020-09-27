package match

import (
	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/services/heartbeat"
)

type MatchService struct {
	Hs   *heartbeat.HeartbeatService
	Repo *db.Repository
}

func (ms *MatchService) MatchOnDemand(sid string, subject string) error {
	// Need to first collect top 10 students, ordered by affinity, send match notifications to each of them (timeout 20 seconds), once a match is found set match Id to a created lesson
	return nil
}
