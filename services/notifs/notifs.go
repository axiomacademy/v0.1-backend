package notifs

import (
	"context"
	"sync"
	// "time"

	log "github.com/sirupsen/logrus"

	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"github.com/solderneer/axiom-backend/db"
	"github.com/solderneer/axiom-backend/graph/model"
)

type NotifService struct {
	logger *log.Logger
	fb     *firebase.App
	fbm    *messaging.Client
	nchans map[string]chan *model.MatchNotification
	nmutex sync.Mutex
}

// Remember to set your GOOGLE_APPLICATION_CREDENTIALS="/home/user/Downloads/service-account-file.json" in the env variables
func (ns *NotifService) Init(logger *log.Logger) {
	ns.logger = logger

	// Initialise firebase
	fb, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		ns.logger.WithFields(log.Fields{
			"service": "notification",
			"error":   err.Error(),
		}).Fatal("Unable to connect to firebase")
	}

	// Create messaging client
	client, err := fb.Messaging(context.Background())
	if err != nil {
		ns.logger.WithFields(log.Fields{
			"service": "notification",
			"error":   err.Error(),
		}).Fatal("Unable to create firebase messaging client")
	}

	ns.fb = fb
	ns.fbm = client
	ns.nchans = map[string]chan *model.MatchNotification{}
	ns.nmutex = sync.Mutex{}

	ns.logger.WithField("service", "notification").Info("Successfully initialised")
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
func (ns *NotifService) SendPushNotification(n db.Notification, token string) error {
	message := &messaging.Message{
		Data: map[string]string{
			"title":    n.Title,
			"subtitle": n.Subtitle,
		},
		Token: token,
	}

	_, err := ns.fbm.Send(context.Background(), message)
	if err != nil {
		ns.logger.WithFields(log.Fields{
			"service":         "notification",
			"notification_id": n.Id,
			"token":           token,
		}).Error("Unable to send push notification")
		return err
	}

	return nil
}
