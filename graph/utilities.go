package graph

import (
	"errors"

	log "github.com/sirupsen/logrus"
)

var InternalServerError = errors.New("Internal Server Error")
var Unauthorised = errors.New("Unauthorised access. Please log in, or switch users to the correct permissions")

func (r *Resolver) sendError(err error, message string) {
	r.Logger.WithFields(log.Fields{
		"service": "resolver",
		"err":     err.Error(),
	}).Error(message)
}
