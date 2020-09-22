package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.
import (
	"github.com/solderneer/axiom-backend/graph/model"
	"sync"
)

type Resolver struct {
	Secret string
	nchans map[string]chan *model.Notification
	mutex  sync.Mutex
}
