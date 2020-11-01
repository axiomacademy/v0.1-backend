package graph

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var (
	InternalServerError = errors.New("Internal Server Error")
	Unauthorised        = errors.New("Unauthorised access. Please log in, or switch users to the correct permissions")
)

// Logs an error with the correct format
func (r *Resolver) sendError(err error, message string) {
	r.Logger.WithFields(log.Fields{
		"service": "resolver",
		"err":     err.Error(),
	}).Error(message)
}
