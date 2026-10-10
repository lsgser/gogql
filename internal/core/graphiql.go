/*
|--------------------------------------------------------------------------
| GraphiQL page
|--------------------------------------------------------------------------
|
| Static HTML template for the default GraphiQL playground. Loads React and
| GraphiQL from unpkg, sets the default query endpoint to the server's
| GraphQLPath, and is injected by playground.go when PlaygroundGraphiQL is selected.
|
| Key func: graphiQLPage(endpoint string) string.
|
*/

package core

import "fmt"

func graphiQLPage(endpoint string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <title>GraphiQL</title>
  <link rel="stylesheet" href="https://unpkg.com/graphiql@3.0.9/graphiql.min.css"/>
</head>
<body style="margin:0;height:100vh;">
  <div id="graphiql" style="height:100%%;">Loading...</div>
  <script crossorigin src="https://unpkg.com/react@18/umd/react.production.min.js"></script>
  <script crossorigin src="https://unpkg.com/react-dom@18/umd/react-dom.production.min.js"></script>
  <script crossorigin src="https://unpkg.com/graphiql@3.0.9/graphiql.min.js"></script>
  <script>
    const fetcher = GraphiQL.createFetcher({ url: %q });
    ReactDOM.createRoot(document.getElementById('graphiql')).render(
      React.createElement(GraphiQL, { fetcher: fetcher })
    );
  </script>
</body>
</html>`, endpoint)
}
