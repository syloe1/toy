package router

import (
	"net/http"

	"url-shortener/internal/handler"
)

func New(urlHandler *handler.URLHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/" && r.Method == http.MethodGet:
			urlHandler.Home(w, r)
		case r.URL.Path == "/api/health":
			urlHandler.Health(w, r)
		case r.URL.Path == "/api/urls" && r.Method == http.MethodPost:
			urlHandler.CreateShortURL(w, r)
		case r.URL.Path == "/api/urls" && r.Method == http.MethodGet:
			urlHandler.ListShortURLs(w, r)
		case len(r.URL.Path) > len("/api/urls/") && r.URL.Path[:len("/api/urls/")] == "/api/urls/":
			urlHandler.GetShortURL(w, r)
		default:
			urlHandler.Redirect(w, r)
		}
	})

	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}
