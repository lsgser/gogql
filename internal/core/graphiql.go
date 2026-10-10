/*
|--------------------------------------------------------------------------
| GraphiQL page
|--------------------------------------------------------------------------
|
| HTML template for the default GraphiQL playground. Scripts and styles load
| from {playgroundPath}/static/ (embedded assets) so the UI works without
| CDN access. Apollo Sandbox still uses its CDN when selected.
|
| Key func: graphiQLPage(graphqlEndpoint, staticAssetPrefix string).
|
*/

package core

import (
	"fmt"
	"strings"
)

func graphiQLPage(graphqlEndpoint, staticAssetPrefix string) string {
	prefix := strings.TrimSuffix(staticAssetPrefix, "/")
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>GraphiQL</title>
  <link rel="stylesheet" href="%s/graphiql.min.css"/>
</head>
<body style="margin:0;height:100vh;">
  <div id="graphiql" style="height:100%%;">Loading GraphiQL…</div>
  <script src="%s/react.production.min.js"></script>
  <script src="%s/react-dom.production.min.js"></script>
  <script src="%s/graphiql.min.js"></script>
  <script>
    const fetcher = GraphiQL.createFetcher({ url: %q });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher: fetcher })
    );
  </script>
</body>
</html>`, prefix, prefix, prefix, prefix, graphqlEndpoint)
}
