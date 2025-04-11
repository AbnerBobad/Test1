package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	//Routes Definitions (GET)
	mux.HandleFunc("GET /", app.loginHandler)          //display home page
	mux.HandleFunc("GET /main", app.mainHandler)       //display main page
	mux.HandleFunc("GET /product", app.productHandler) //display product page
	mux.HandleFunc("GET /view", app.viewHandler)       //display view page

	//Submit form
	mux.HandleFunc("POST /product", app.createProduct) //submit product form

	//Retrieve Data GETS

	//Routes Definitions (POST)

	return app.loggingMiddleware(mux)
}
