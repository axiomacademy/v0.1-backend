package notifs

import (
	"github.com/solderneer/axiom-backend/graph/model"
	"sync"
)

type NotifService struct {
	nchans map[string]chan *model.Notification
	mutex  sync.Mutex
}

func (ns *NotifService) SendNotification() {

}
