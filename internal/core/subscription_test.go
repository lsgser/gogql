/*
|--------------------------------------------------------------------------
| Subscription tests
|--------------------------------------------------------------------------
|
| Exercises Application.Subscribe with a module that uses SubscriptionResolvers
| (struct methods), ensuring the subscription root is wired separately from
| ResolverMap.
|
*/

package core_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/lsgser/gogql"
)

func TestSubscription(t *testing.T) {
	mod := gogql.MustModule(gogql.ModuleConfig{
		ID: "sub",
		TypeDefs: `
			type Query { ok: String! }
			type Event { msg: String! }
			type Subscription { events: Event! }
		`,
		Resolvers: gogql.NewResolverMap().
			Query("ok", func(_ context.Context) (string, error) { return "ok", nil }),
		SubscriptionResolvers: &subscriptionRoot{},
	})

	app := gogql.MustApplication(gogql.ApplicationConfig{Modules: []*gogql.Module{mod}})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	ch, err := app.Subscribe(ctx, `subscription { events { msg } }`, "", nil)
	if err != nil {
		t.Fatal(err)
	}

	select {
	case v := <-ch:
		raw, err := json.Marshal(v)
		if err != nil {
			t.Fatal(err)
		}
		var resp struct {
			Data struct {
				Events struct {
					Msg string `json:"msg"`
				} `json:"events"`
			} `json:"data"`
			Errors []any `json:"errors"`
		}
		if err := json.Unmarshal(raw, &resp); err != nil {
			t.Fatal(err)
		}
		if len(resp.Errors) > 0 {
			t.Fatalf("errors: %+v", resp.Errors)
		}
		if resp.Data.Events.Msg != "hi" {
			t.Fatalf("msg: %q", resp.Data.Events.Msg)
		}
	case <-ctx.Done():
		t.Fatal("timeout waiting for subscription event")
	}
}

type subscriptionRoot struct{}

func (subscriptionRoot) Events(ctx context.Context) (chan *event, error) {
	ch := make(chan *event, 1)
	ch <- &event{msg: "hi"}
	close(ch)
	return ch, nil
}

type event struct {
	msg string
}

func (e *event) Msg() string {
	if e == nil {
		return ""
	}
	return e.msg
}
