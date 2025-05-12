package main

import (
	"net/http"
	"time"
)

func (app *application) serve() error {

	srv := &http.Server{
		Addr:         *app.addr,
		Handler:      app.routes(),
		TLSConfig:    app.tlsConfig, //new
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	app.logger.Info("Starting server", "addr", srv.Addr)
	return srv.ListenAndServeTLS("./tls/cert.pem", "./tls/key.pem") //New changes
}
