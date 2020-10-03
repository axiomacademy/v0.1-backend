package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/services/notifs"
	"github.com/solderneer/axiom-backend/services/video"
)

type Resolver struct {
	Secret string
	Repo   *db.Repository
	Ns     *notifs.NotifService
	Video	 *video.VideoClient
}
