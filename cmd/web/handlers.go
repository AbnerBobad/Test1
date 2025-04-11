package main

import (
	"net/http"
)

// LOGIN START
// LoginHandler is a handler that renders the home page - index.tmpl
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Welcome to StockTrack, your stock tracking website!"
	err := app.render(w, http.StatusOK, "index.html", data)
	if err != nil {
		app.logger.Error("failed to render the Login page", "template", "index.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// MAIN PAGE START
// mainHandler is a handler that renders the main page - main.tmpl
func (app *application) mainHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Welcome to the Main Page"
	err := app.render(w, http.StatusOK, "main.html", data)
	if err != nil {
		app.logger.Error("failed to render the Main Page", "template", "main.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PRODUCT START
// productHandler is a handler that render the product page - product.tmpl
func (app *application) productHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Insert your product here"
	err := app.render(w, http.StatusOK, "product.html", data)
	if err != nil {
		app.logger.Error("failed to render the Product Page", "template", "product.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// VIEW START
// viewHandler is a handler that renders the view page - view.tmpl
func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "View your product's here"
	err := app.render(w, http.StatusOK, "view.html", data)
	if err != nil {
		app.logger.Error("failed to render the view Page", "template", "view.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
