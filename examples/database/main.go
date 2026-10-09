/*
|--------------------------------------------------------------------------
| Database example — entrypoint
|--------------------------------------------------------------------------
|
| Opens SQLite, builds MustApplication with loaders, JWT ContextFunc,
| and serves GraphQL + playground.
|
*/

package main

import (
	"log"
	"time"

	"github.com/lsgser/gogql"
	"github.com/lsgser/gogql/examples/database/auth"
	"github.com/lsgser/gogql/examples/database/modules"
	"github.com/lsgser/gogql/examples/database/store"
)

func main() {
	db, err := store.OpenSQLite("file:example.db?cache=shared&mode=rwc")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := store.NewUserRepository(db)

	token, err := auth.MintDemoToken("1", 24*time.Hour)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Demo JWT (user id 1):", token)
	log.Println("Use header: Authorization: Bearer", token)

	app := gogql.MustApplication(gogql.ApplicationConfig{
		Modules:  modules.All(repo),
		Loaders:  modules.LoaderFactories(repo),
		Security: gogql.SecurityConfig{MaxDepth: 10},
	})

	server := gogql.NewServer(app, gogql.ServerConfig{
		ContextFunc: auth.ContextFunc([]byte(auth.DemoSecret)),
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
