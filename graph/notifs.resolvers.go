package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/solderneer/axiom-backend/graph/generated"
	"github.com/solderneer/axiom-backend/graph/model"
)

func (r *subscriptionResolver) SubscribeNotifications(ctx context.Context, user string) (<-chan *model.Notification, error) {
	// Creating the channel
	nchan := make(chan *model.Notification, 1)
	r.ns.mutex.Lock()
	r.ns.nchans[user] = nchan
	r.ns.mutex.Unlock()

	// Delete channel when done
	go func() {
		<-ctx.Done()
		r.ns.mutex.Lock()
		delete(r.ns.nchans, user)
		r.ns.mutex.Unlock()
	}()

	return nchan, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
