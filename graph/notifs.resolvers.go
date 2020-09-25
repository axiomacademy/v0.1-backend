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
	r.Ns.Nmutex.Lock()
	r.Ns.Nchans[user] = nchan
	r.Ns.Nmutex.Unlock()

	// Delete channel when done
	go func() {
		<-ctx.Done()
		r.Ns.Nmutex.Lock()
		delete(r.Ns.Nchans, user)
		r.Ns.Nmutex.Unlock()
	}()

	return nchan, nil
}

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type subscriptionResolver struct{ *Resolver }
