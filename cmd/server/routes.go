// package main

// import (
// 	"lxr-d/internal/handlers"

// 	"github.com/go-chi/chi/v5"
// )

// func NewRouter(h *handlers.Handler) *chi.Mux {

// 	r := chi.NewRouter()

// 	r.Get("/ping", h.PingHanlder)

// 	r.Post("/create", h.CreateHandler)
// 	r.Post("/start", h.StartHandler)
// 	r.Get("/stop", h.StopHandler)
// 	r.Get("/exec", h.ExecHandler)
// 	r.Delete("/kill", h.KillHanlder)
// 	r.Post("/pull_image", h.PullImageHandler)

// 	r.Get("/ps", h.PsHandler)
// 	r.Get("/ps/all", h.PsAllHandler)

// 	return r
// }

package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"lxr-d/internal/handlers"

	"github.com/go-chi/chi/v5"
)

// NewRouter builds the single router shared by the Unix-socket API server and
// the TCP web server. The existing daemon routes are kept intact for CLI/API
// compatibility, while /api/* exposes the same handlers to the web console.
func NewRouter(h *handlers.Handler) *chi.Mux {
	r := chi.NewRouter()

	registerAPIRoutes(r, h, "")
	registerAPIRoutes(r, h, "/api")

	// Everything that is not an API/legacy daemon route is handled by the
	// embedded web console. API 404s therefore remain API 404s because the
	// explicit routes above take precedence over this catch-all.
	r.Handle("/*", webHandler())

	return r
}

func registerAPIRoutes(r chi.Router, h *handlers.Handler, prefix string) {
	// Health
	r.Get(prefix+"/ping", h.PingHanlder)

	// Container lifecycle
	r.Post(prefix+"/create", h.CreateHandler)
	r.Post(prefix+"/start", h.StartHandler)
	r.Get(prefix+"/stop", h.StopHandler)
	r.Get(prefix+"/exec", h.ExecHandler)
	r.Delete(prefix+"/kill", h.KillHanlder)

	// Images
	r.Post(prefix+"/pull_image", h.PullImageHandler)

	// Container listing
	r.Get(prefix+"/ps", h.PsHandler)
	r.Get(prefix+"/ps/all", h.PsAllHandler)
}

// webHandler serves the static console and falls back to index.html for SPA
// routes such as /containers/foo and /containers/foo/terminal.
func webHandler() http.Handler {
	root := os.Getenv("LXR_WEB_DIR")
	if root == "" {
		// The daemon is normally started from cmd/server (see Makefile).
		root = "../../web"
	}

	root, err := filepath.Abs(root)
	if err != nil {
		// Keep the handler alive; requests will return 500 below if the path
		// cannot be resolved.
		return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, "web root is unavailable", http.StatusInternalServerError)
		})
	}

	fileServer := http.FileServer(http.Dir(root))
	indexPath := filepath.Join(root, "index.html")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Never let the web fallback swallow an API request. This also makes
		// unknown API paths return a normal JSON-ish HTTP 404 instead of HTML.
		if strings.HasPrefix(r.URL.Path, "/api/") || r.URL.Path == "/api" {
			http.NotFound(w, r)
			return
		}

		path := filepath.Clean(r.URL.Path)
		if path == "/" {
			http.ServeFile(w, r, indexPath)
			return
		}

		// Check for a real file before applying the SPA fallback.
		requested := filepath.Join(root, filepath.FromSlash(strings.TrimPrefix(path, "/")))
		if info, statErr := os.Stat(requested); statErr == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Browser history routes belong to the vanilla JS SPA.
		http.ServeFile(w, r, indexPath)
	})
}
