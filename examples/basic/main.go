/*
|--------------------------------------------------------------------------
| Basic example — entrypoint
|--------------------------------------------------------------------------
|
| Wires modules.All() into MustApplication and starts the HTTP server
| with GraphiQL playground.
|
*/

package main

import (
	"log"

	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/basic/modules"
)

func main() {
	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules: modules.All(),
	})

	server := gogql.NewServer(app, gogql.ServerConfig{
		Playground: &gogql.PlaygroundConfig{
			Enabled: true,
			Path:    "/playground",
			UI:      gogql.PlaygroundGraphiQL,
		},
	})

	log.Println("GraphQL:    http://localhost:8080/graphql")
	log.Println("Playground: http://localhost:8080/playground")
	log.Fatal(server.ListenAndServe(":8080"))
}
