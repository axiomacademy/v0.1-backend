package db

import (
	"context"
	"time"

	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Notification struct {
	Id       string
	Title    string
	Subtitle string
	Image    string
	Read     bool
	Created  time.Time
}

func (r *Repository) CreateNotification() (Notification, error) {
	return Notification{}, nil
}

func (r *Repository) UpdateNotification(n Notification) error {
	return nil
}

func (r *Repository) GetUserNotifications(uid string, limit int, offset int) ([]Notification, error) {
	return nil, nil
}
