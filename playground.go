package gogql

import "fmt"

func renderPlayground(settings playgroundSettings) (string, error) {
	switch settings.ui {
	case PlaygroundApolloSandbox:
		return apolloSandboxPage(settings.graphqlPath), nil
	case PlaygroundGraphiQL:
		return graphiQLPage(settings.graphqlPath), nil
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
