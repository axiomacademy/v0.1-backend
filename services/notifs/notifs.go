package notifs

import (
	"context"
	"sync"
	"time"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Notification struct {
	Image    string
	Title    string
	Subtitle string
	Read     bool
	Created  time.Time
}

type NotifService struct {
	fb     *firebase.App
	fbm    *messaging.Client
	nchans map[string]chan *model.MatchNotification
	nmutex sync.Mutex
}

// Remember to set your GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json" in the env variables
func (ns *NotifService) Init() error {
	// Initialise firebase
	fb, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return err
	}

	// Create messaging client
	client, err := fb.Messaging(context.Background())
	if err != nil {
		return err
	}

	ns.fb = fb
	ns.fbm = client
	ns.nchans = map[string]chan *model.MatchNotification{}
	ns.nmutex = sync.Mutex{}

	return nil
}

func (ns *NotifService) SendMatchNotification(n model.MatchNotification, uid string) {
	ns.nchans[uid] <- &n
}

func (ns *NotifService) CreateUserMatchChannel(user string) *chan *model.MatchNotification {
	nchan := make(chan *model.MatchNotification, 1)

	ns.nmutex.Lock()
	ns.nchans[user] = nchan
	ns.nmutex.Unlock()

	return &nchan
}

func (ns *NotifService) DeleteUserMatchChannel(user string) {
	ns.nmutex.Lock()
	delete(ns.nchans, user)
	ns.nmutex.Unlock()
}

// Takes a notification struct and the registration token of the user. Push notifications omit the image of the notification and time
func (ns *NotifService) SendPushNotification(n Notification, token string) error {
	message := &messaging.Message{
		Data: map[string]string{
			"title":    n.Title,
			"subtitle": n.Subtitle,
		},
		Token: token,
	}

	_, err := ns.fbm.Send(context.Background(), message)
	if err != nil {
		return err
	}

	return nil
}
