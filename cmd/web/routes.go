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
	// mux.HandleFunc("GET /", app.mainHandler)           //display home page
	// mux.HandleFunc("GET /product", app.productHandler) //display product page
	// mux.HandleFunc("GET /view", app.viewHandler)       //display view page
	// mux.HandleFunc("GET /login", app.loginHandler)     //display login page

	// //Submit form
	// mux.HandleFunc("POST /product", app.createProduct) //submit product form

	// mux.HandleFunc("GET /product/edit", app.editProductForm)  //get the changed data
	// mux.HandleFunc("POST /product/update", app.updateProduct) //update data
	// mux.HandleFunc("POST /product/delete", app.deleteProduct) //delete data

	// mux.HandleFunc("GET /search", app.searchProducts) //search for data

	//MAIN
	mux.Handle("GET /", http.HandlerFunc(app.mainHandler))
	mux.Handle("GET /product", http.HandlerFunc(app.productHandler))

	//LOGIN
	mux.Handle("GET /login", app.session.Enable(http.HandlerFunc(app.loginHandler)))
	mux.Handle("GET /user/signup", app.session.Enable(http.HandlerFunc(app.signupUserForm)))
	mux.Handle("POST /user/signup", app.session.Enable(http.HandlerFunc(app.signupUser)))

	//FUNCTIONALITY
	mux.Handle("GET /view", app.session.Enable(http.HandlerFunc(app.viewHandler)))
	mux.Handle("POST /product", app.session.Enable(http.HandlerFunc(app.createProduct)))
	mux.Handle("GET /product/edit", app.session.Enable(http.HandlerFunc(app.editProductForm)))
	mux.Handle("POST /product/update", app.session.Enable(http.HandlerFunc(app.updateProduct)))
	mux.Handle("POST /product/delete", app.session.Enable(http.HandlerFunc(app.deleteProduct)))
	mux.Handle("GET /search", app.session.Enable(http.HandlerFunc(app.searchProducts)))

	return app.loggingMiddleware(mux)
}
