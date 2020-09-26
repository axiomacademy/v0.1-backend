// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type User interface {
	IsUser()
}

type Heartbeat struct {
	Status   HeartbeatStatus `json:"status"`
	LastSeen int             `json:"lastSeen"`
}

type Lesson struct {
	ID       string   `json:"id"`
	Subject  string   `json:"subject"`
	Summary  string   `json:"summary"`
	Tutor    *Tutor   `json:"tutor"`
	Student  *Student `json:"student"`
	Duration int      `json:"duration"`
	Date     string   `json:"date"`
	Chat     string   `json:"chat"`
}

type LoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type NewStudent struct {
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProfilePic string `json:"profilePic"`
}

type NewTutor struct {
	Username   string   `json:"username"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	Password   string   `json:"password"`
	ProfilePic string   `json:"profilePic"`
	HourlyRate int      `json:"hourlyRate"`
	Bio        string   `json:"bio"`
	Education  []string `json:"education"`
	Subjects   []string `json:"subjects"`
}

type Notification struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type Student struct {
	ID         string `json:"id"`
	Username   string `json:"username"`
	FirstName  string `json:"firstName"`
	LastName   string `json:"lastName"`
	Email      string `json:"email"`
	ProfilePic string `json:"profilePic"`
}

func (Student) IsUser() {}

type Tutor struct {
	ID         string   `json:"id"`
	Username   string   `json:"username"`
	FirstName  string   `json:"firstName"`
	LastName   string   `json:"lastName"`
	Email      string   `json:"email"`
	ProfilePic string   `json:"profilePic"`
	HourlyRate int      `json:"hourlyRate"`
	Bio        string   `json:"bio"`
	Rating     int      `json:"rating"`
	Education  []string `json:"education"`
	Subjects   []string `json:"subjects"`
}

func (Tutor) IsUser() {}

type HeartbeatStatus string

const (
	HeartbeatStatusOnline  HeartbeatStatus = "ONLINE"
	HeartbeatStatusActive  HeartbeatStatus = "ACTIVE"
	HeartbeatStatusOffline HeartbeatStatus = "OFFLINE"
)

var AllHeartbeatStatus = []HeartbeatStatus{
	HeartbeatStatusOnline,
	HeartbeatStatusActive,
	HeartbeatStatusOffline,
}

func (e HeartbeatStatus) IsValid() bool {
	switch e {
	case HeartbeatStatusOnline, HeartbeatStatusActive, HeartbeatStatusOffline:
		return true
	}
	return false
}

func (e HeartbeatStatus) String() string {
	return string(e)
}

func (e *HeartbeatStatus) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = HeartbeatStatus(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid HeartbeatStatus", str)
	}
	return nil
}

func (e HeartbeatStatus) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
