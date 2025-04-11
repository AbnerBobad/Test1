package main

import (
	"net/http"
	"strconv"

	"github.com/AbnerBobad/final_project/internal/data"
	"github.com/AbnerBobad/final_project/internal/validator"
)

// LOGIN START
// LoginHandler is a handler that renders the home page - index.tmpl
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Login Page"
	data.FileInfo = "Please login to continue."
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
	data.HeaderText = "Welcome to StockTrack"
	data.FileInfo = "Manage your inventory efficiently and stay updated on stock levels."
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
	data.HeaderText = "Add New Products"
	data.FileInfo = "Please fill in the product details below."
	err := app.render(w, http.StatusOK, "product.html", data)
	if err != nil {
		app.logger.Error("failed to render the Product Page", "template", "product.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// product form creation
func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	//parsed data form
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse products from data", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//data members parsed
	productName := r.PostForm.Get("product_name")
	productQuantityStr := r.PostForm.Get("product_quantity")
	productPriceStr := r.PostForm.Get("product_price")
	productDescription := r.PostForm.Get("product_description")

	productQuantity, err := strconv.ParseInt(productQuantityStr, 10, 64)
	if err != nil {
		app.logger.Error("invalid product quantity", "error", err)
		http.Error(w, "Invalid product quantity", http.StatusBadRequest)
		return
	}

	productPrice, err := strconv.ParseFloat(productPriceStr, 64)
	if err != nil {
		app.logger.Error("invalid product price", "error", err)
		http.Error(w, "Invalid product price", http.StatusBadRequest)
		return
	}

	//Instance for data members
	product := &data.Product{
		PName:        productName,
		PQuantity:    productQuantity,
		PPrice:       productPrice,
		PDescription: productDescription,
	}
	//Data Validator
	v := validator.NewValidator()
	data.ValidateProduct(v, product)
	//check for validation
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "StockTrack"
		data.HeaderText = "Add New Products"
		data.FileInfo = "Please fill in the product details below."
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"product_name":        productName,
			"product_quantity":    productQuantityStr,
			"product_price":       productPriceStr,
			"product_description": productDescription,
		}
		err := app.render(w, http.StatusOK, "product.html", data)
		if err != nil {
			app.logger.Error("failed to render the Product Page", "template", "product.html", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	//error checker
	err = app.products.Insert(product)
	if err != nil {
		app.logger.Error("failed to insert product into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/main", http.StatusSeeOther)
}

// VIEW START
// viewHandler is a handler that renders the view page - view.tmpl
func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Current Inventory"
	data.FileInfo = "Here are the products in your inventory."
	err := app.render(w, http.StatusOK, "view.html", data)
	if err != nil {
		app.logger.Error("failed to render the view Page", "template", "view.html", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
