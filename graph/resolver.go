// Package graph contains all the GraphQL resolvers
package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	log "github.com/sirupsen/logrus"

	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/services/chat"
	"github.com/solderneer/axiom-backend/services/match"
	"github.com/solderneer/axiom-backend/services/notifs"
	"github.com/solderneer/axiom-backend/services/video"
)

type Resolver struct {
	Secret string
	Logger *log.Logger
	Repo   *db.Repository
	Ns     *notifs.NotifService
	Cs     *chat.Chat
	Video  *video.VideoClient
	Ms     *match.MatchService
}
