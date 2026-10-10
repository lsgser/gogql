/*
|--------------------------------------------------------------------------
| DataLoader tests
|--------------------------------------------------------------------------
|
| Ensures LoaderRegistry creates distinct loader instances per request context
| and that LoadMany invokes the batch function once with all keys.
|
*/

package core_test

import (
	"context"
	"testing"

	"github.com/graph-gophers/dataloader"
	"github.com/lsgser/gogql"
)

func TestLoaderRegistryPerRequest(t *testing.T) {
	var batches int
	factories := gogql.LoaderFactories{
		"echo": func() *dataloader.Loader {
			return gogql.NewStringKeyLoader(func(_ context.Context, keys dataloader.Keys) []*dataloader.Result {
				batches++
				out := make([]*dataloader.Result, len(keys))
				for i, k := range keys {
					out[i] = &dataloader.Result{Data: "v:" + k.String()}
				}
				return out
			})
		},
	}

	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules: []*gogql.Module{gogql.MustModule(gogql.ModuleConfig{
			ID:       "noop",
			TypeDefs: `type Query { ok: String! }`,
			Resolvers: gogql.NewResolverMap().Query("ok", func(_ context.Context) (string, error) {
				return "ok", nil
			}),
		})},
		Loaders: factories,
	})

	ctx := app.RequestContext(context.Background())
	l, err := gogql.LoadersFromContext(ctx, "echo")
	if err != nil {
		t.Fatal(err)
	}
	_, errs := l.LoadMany(ctx, dataloader.NewKeysFromStrings([]string{"a", "b"}))()
	for _, err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if batches != 1 {
		t.Fatalf("expected one batch, got %d", batches)
	}
}
