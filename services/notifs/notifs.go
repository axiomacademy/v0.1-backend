package notifs

import (
	"github.com/solderneer/axiom-backend/graph/model"
	"sync"
)

type NotifService struct {
	Nchans map[string]chan *model.Notification
	Nmutex sync.Mutex
}

func (ns *NotifService) SendNotification(n model.Notification, uid string) {
	ns.Nchans[uid] <- &n
}
