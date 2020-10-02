package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/services/match"
	"github.com/solderneer/axiom-backend/services/notifs"
)

type Resolver struct {
	Secret string
	Repo   *db.Repository
	Ns     *notifs.NotifService
	Ms     *match.MatchService
}
