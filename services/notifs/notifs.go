package notifs

import (
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
	ns.nchans[uid] <- &n
}

func (ns *NotifService) CreateUserChannel(user string) *chan *model.Notification {
	nchan := make(chan *model.Notification, 1)

	ns.nmutex.Lock()
	ns.nchans[user] = nchan
	ns.nmutex.Unlock()

	return &nchan
}

func (ns *NotifService) DeleteUserChannel(user string) {
	ns.nmutex.Lock()
	delete(ns.nchans, user)
	ns.nmutex.Unlock()
}

func (ns *NotifService) SendPushNotification(n model.Notification, uid string) {
	return
}
