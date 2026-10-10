/*
|--------------------------------------------------------------------------
| Playground UI
|--------------------------------------------------------------------------
|
| Builds the HTML page returned at the playground route. renderPlayground
| switches on PlaygroundUI (GraphiQL vs Apollo Sandbox), embedding CDN
| scripts and pointing the UI at the configured GraphQL HTTP path.
|
| Used only from server.go handlePlayground; not a public API for apps.
|
| Key funcs: renderPlayground, apolloSandboxPage (GraphiQL HTML lives in graphiql.go).
|
*/

package core

import (
	"fmt"
	"strings"
)

func renderPlayground(settings playgroundSettings) (string, error) {
	switch settings.ui {
	case PlaygroundApolloSandbox:
		return apolloSandboxPage(settings.graphqlPath), nil
	case PlaygroundGraphiQL:
		staticPrefix := strings.TrimSuffix(settings.path, "/") + "/static"
		return graphiQLPage(settings.graphqlPath, staticPrefix), nil
	default:
		return "", fmt.Errorf("gogql: unknown playground UI %q", settings.ui)
	}
}

func apolloSandboxPage(endpoint string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>Apollo Sandbox</title>
  <style>html,body,#sandbox{height:100%%;margin:0;}</style>
</head>
<body>
  <div id="sandbox"></div>
  <script src="https://embeddable-sandbox.cdn.apollographql.com/_latest/embeddable-sandbox.umd.production.min.js"></script>
  <script>
    new window.EmbeddedSandbox({
      target: '#sandbox',
      initialEndpoint: window.location.origin + %q,
    });
  </script>
</body>
</html>`, endpoint)
}
