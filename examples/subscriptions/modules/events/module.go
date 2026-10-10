/*
|--------------------------------------------------------------------------
| Subscriptions example — events module
|--------------------------------------------------------------------------
|
| Demonstrates SubscriptionResolvers (struct method Ticks returning a channel)
| alongside ResolverMap for Query—subscriptions cannot use ResolverMap.Subscription.
| WebSocket context uses the same ContextFunc as HTTP when configured on Server.
|
*/

package events

import (
	"context"
	"time"

	"github.com/lsgser/gogql"
)

// Module defines query + subscription SDL and gogql resolvers.
func Module() *gogql.Module {
	return gogql.MustModule(gogql.ModuleConfig{
		ID: "events",
		TypeDefs: `
			type Query {
				ok: String!
			}
			type TickEvent {
				at: String!
				count: Int!
			}
			type Subscription {
				ticks: TickEvent!
			}
		`,
		Resolvers: gogql.NewResolverMap().
			Query("ok", func(_ context.Context) (string, error) { return "ok", nil }),
		SubscriptionResolvers: &subscriptionRoot{},
	})
}

type subscriptionRoot struct{}

func (subscriptionRoot) Ticks(ctx context.Context) (chan *TickEvent, error) {
	ch := make(chan *TickEvent, 1)
	go func() {
		defer close(ch)
		for i := int32(1); i <= 3; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- &TickEvent{at: time.Now().Format(time.RFC3339), count: i}:
				time.Sleep(time.Second)
			}
		}
	}()
	return ch, nil
}

// TickEvent is the subscription payload type.
type TickEvent struct {
	at    string
	count int32
}

func (e *TickEvent) At() string   { return e.at }
func (e *TickEvent) Count() int32 { return e.count }
