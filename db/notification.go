package db

import (
	"context"
	"strings"
	"time"

	"github.com/jackc/pgtype"
	"github.com/pborman/uuid"
	"github.com/solderneer/axiom-backend/graph/model"
)

type Notification struct {
	Id       string
	Tutor    string
	Student  string
	Title    string
	Subtitle string
	Image    string
	Read     bool
	Created  time.Time
}

// Convert a db.Notification to a model.Notification
func (r *Repository) ToNotificationModel(n Notification) model.Notification {
	return model.Notification{ID: n.Id, Title: n.Title, Subtitle: n.Subtitle, Image: n.Image, Created: n.Created}
}

// Create a new notification and commits it to the database
// Takes the user's ID (either Tutor, Student), title, subtitle and image
func (r *Repository) CreateNotification(uid string, title string, subtitle string, image string) (Notification, error) {
	var n Notification

	n.Id = uuid.New()
	n.Title = title
	n.Subtitle = subtitle
	n.Image = image
	n.Read = false
	n.Created = time.Now()

	// Parse uid
	idSplit := strings.Split(uid, ":")
	if idSplit[0] == "s" {
		n.Student = uid
	} else if idSplit[0] == "t" {
		n.Tutor = uid
	}

	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return n, err
	}

	defer tx.Rollback(context.Background())

	sql := `INSERT INTO notifications (id, tutor, student, title, subtitle, image, read, created) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err = tx.Exec(context.Background(), sql, n.Id, n.Tutor, n.Student, n.Title, n.Subtitle, n.Image, n.Read, n.Created)

	if err != nil {
		return n, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return n, err
	}

	return n, nil
}

// Update notification based on an existing notification struct, only can update read status
func (r *Repository) UpdateNotification(n Notification) error {
	tx, err := r.dbPool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())
	sql := `UPDATE notifications SET read = $2 WHERE id = $1`
	_, err = tx.Exec(context.Background(), sql, n.Id, n.Read)

	if err != nil {
		return err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

// Get notification by notification UUID
func (r *Repository) GetNotificationById(nid string) (Notification, error) {
	sql := `SELECT id, tutor, student, title, subtitle, image, read, created FROM notifications WHERE id = $1`

	var n Notification

	// To handle possible null values
	var s pgtype.Varchar
	var t pgtype.Varchar

	if err := r.dbPool.QueryRow(context.Background(), sql, nid).Scan(&n.Id, &t, &s, &n.Title, &n.Image, &n.Read, &n.Created); err != nil {
		return n, err
	}

	// Handling possible null values
	s.AssignTo(&n.Student)
	t.AssignTo(&n.Tutor)
	return n, nil
}

// Get all the notifications associated with a user. Paginated by time period
func (r *Repository) GetUserNotifications(uid string, startTime time.Time, endTime time.Time) ([]Notification, error) {
	var sql string

	idSplit := strings.Split(uid, ":")
	if idSplit[0] == "s" {
		sql = `SELECT id, tutor, student, title, subtitle, image, read, created FROM notifications WHERE student = $1 AND created > $2 AND created < $3`
	} else if idSplit[0] == "t" {
		sql = `SELECT id, tutor, student, title, subtitle, image, read, created FROM notifications WHERE tutor = $1 AND created > $2 AND created < $3`
	}

	var notifications []Notification

	var st pgtype.Timestamptz
	var et pgtype.Timestamptz

	st.Time = startTime
	st.Status = pgtype.Present

	et.Time = endTime
	et.Status = pgtype.Present

	rows, err := r.dbPool.Query(context.Background(), sql, uid, st, et)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		var n Notification

		// To handle possible null values
		var s pgtype.Varchar
		var t pgtype.Varchar

		if err := rows.Scan(&n.Id, &t, &s, &n.Tutor, &n.Subtitle, &n.Image, &n.Read, &n.Created); err != nil {
			return nil, err
		}

		// Handling possible null values
		s.AssignTo(&n.Student)
		t.AssignTo(&n.Tutor)
		notifications = append(notifications, n)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return notifications, nil
}
