package notifs

import (
	"http"
	"sync"

	"github.com/solderneer/axiom-backend/graph/model"
)

type NotifService struct {
	nchans map[string]chan *model.Notification
	nmutex sync.Mutex
}

func (ns *NotifService) Init() {
	ns.nchans = map[string]chan *model.Notification{}
	ns.nmutex = sync.Mutex{}
}

func (ns *NotifService) SendNotification(n model.Notification, uid string) {
	ns.nmutex.Lock()
	ns.nchans[uid] <- &n
	ns.nmutex.Unlock()
}

func (ns *NotifService) SendPushNotification(n model.Notification, uid string) {
	return
}
