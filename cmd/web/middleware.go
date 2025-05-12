package main

import (
	"net/http"
)

// Middleware for logging
func (app *application) loggingMiddleware(next http.Handler) http.Handler {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Info("received request", "ip", ip, "protocol", proto, "method", method, "uri", uri)
		next.ServeHTTP(w, r)
		app.logger.Info("Request processed")
	})
	return fn

}

// required middleware for redirection to login is no
func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !app.IsAuthenticated(r) {
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		w.Header().Add("Cache Control", "no-store")
		next.ServeHTTP(w, r)

	})
}
