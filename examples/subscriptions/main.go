/*
|--------------------------------------------------------------------------
| Subscriptions example — entrypoint
|--------------------------------------------------------------------------
|
| Serves HTTP GraphQL and graphql-transport-ws on the same path. Connect a
| client with subprotocol graphql-transport-ws to exercise Subscription fields.
| Run: go run ./examples/subscriptions
|
*/

package main

import (
	"log"

	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/subscriptions/modules"
)

func main() {
	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules:  modules.All(),
		Security: gogql.SecurityConfig{MaxDepth: 10},
	})

	server := gogql.NewServer(app, gogql.ServerConfig{
		Playground: &gogql.PlaygroundConfig{
			Enabled: true,
			Path:    "/playground",
			UI:      gogql.PlaygroundGraphiQL,
		},
	})

	log.Println("GraphQL + subscriptions: ws://localhost:8080/graphql (graphql-transport-ws)")
	log.Println("Playground: http://localhost:8080/playground")
	log.Fatal(server.ListenAndServe(":8080"))
}
