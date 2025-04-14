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
	mux.HandleFunc("GET /", app.mainHandler)           //display home page
	mux.HandleFunc("GET /product", app.productHandler) //display product page
	mux.HandleFunc("GET /view", app.viewHandler)       //display view page

	//Submit form
	mux.HandleFunc("POST /product", app.createProduct) //submit product form

	mux.HandleFunc("GET /product/edit", app.editProductForm)
	mux.HandleFunc("POST /product/update", app.updateProduct)
	mux.HandleFunc("POST /product/delete", app.deleteProduct)

	return app.loggingMiddleware(mux)
}
