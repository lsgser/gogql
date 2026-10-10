/*
|--------------------------------------------------------------------------
| Playground static assets
|--------------------------------------------------------------------------
|
| GraphiQL, React, and ReactDOM bundles embedded for offline/local playground
| use (no unpkg.com required). Served under {PlaygroundPath}/static/.
|
*/

package playground

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed static/*
var staticFiles embed.FS

// StaticFS returns the embedded GraphiQL asset filesystem.
func StaticFS() (fs.FS, error) {
	return fs.Sub(staticFiles, "static")
}

// StaticHandler serves embedded playground assets.
func StaticHandler() (http.Handler, error) {
	sub, err := StaticFS()
	if err != nil {
		return nil, err
	}
	return http.FileServer(http.FS(sub)), nil
}
