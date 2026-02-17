package pathhandler

import (
	"net/http"
	"strings"
)

func PathHandlerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path != "/" && !strings.HasSuffix(path, "/") {
			r.URL.Path = path + "/"
		}
		next.ServeHTTP(w, r)
	})
}
