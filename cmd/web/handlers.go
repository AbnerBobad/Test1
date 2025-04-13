package main

import (
	"net/http"
	"strconv"

	"github.com/AbnerBobad/final_project/internal/data"
)

// LOGIN START
// LoginHandler is a handler that renders the home page - index.tmpl
// func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
// 	data := NewTemplateData()
// 	data.Title = "StockTrack"
// 	data.HeaderText = "Login Page"
// 	data.FileInfo = "Please login to continue."
// 	err := app.render(w, http.StatusOK, "index.tmpl", data)
// 	if err != nil {
// 		app.logger.Error("failed to render the Login page", "template", "index.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		return
// 	}
// }

// MAIN PAGE START
// mainHandler is a handler that renders the main page - main.tmpl
func (app *application) mainHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Welcome to StockTrack"
	data.FileInfo = "Manage your inventory efficiently and stay updated on stock levels."
	err := app.render(w, http.StatusOK, "main.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Main Page", "template", "main.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
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
	err := app.render(w, http.StatusOK, "product.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Product Page", "template", "product.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
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
	//check if the product name, quantity, price and description are empty
	if productName == "" || productQuantityStr == "" || productPriceStr == "" || productDescription == "" {
		data := NewTemplateData()
		data.Title = "StockTrack"
		data.HeaderText = "Add New Products"
		data.FileInfo = "Please fill in the product details below."
		data.FormErrors = map[string]string{}
		//check if the product name is valid
		if productName == "" {
			data.FormErrors["product_name"] = "Product name is required"
		} else if len(productName) < 1 || len(productName) > 100 {
			data.FormErrors["product_name"] = "Product Name must be less than 100 characters"
		}
		// Check if the product name contains any numbers, CANT BE IMPLEMENTED
		// because the product name can contain numbers
		// for _, char := range productName {
		// 	if char >= '0' && char <= '9' {
		// 		data.FormErrors["product_name"] = "Product name must not contain numbers"
		// 		break
		// 	}
		// }

		//check if the product quantity is valid
		productQuantity, err := strconv.ParseInt(productQuantityStr, 10, 64)
		//check if the product quantity is valid
		if productQuantityStr == "" {
			data.FormErrors["product_quantity"] = "Product quantity is required"
		} else if err != nil || productQuantity <= 0 {
			data.FormErrors["product_quantity"] = "Product Quantity must be a positive number"
		}
		//parse the product price
		productPrice, err := strconv.ParseFloat(productPriceStr, 64)
		//check if the product price is valid
		if productPriceStr == "" {
			data.FormErrors["product_price"] = "Product Price is required"
		} else if err != nil || productPrice <= 0.0 {
			data.FormErrors["product_price"] = "Product Price must be a positive number"
		}
		if productDescription == "" {
			data.FormErrors["product_description"] = "Product description is required"
		}
		// if productDescription == "" {
		// 	productDescription = "none"
		// }
		data.FormData = map[string]string{
			"product_name":        productName,
			"product_quantity":    productQuantityStr,
			"product_price":       productPriceStr,
			"product_description": productDescription,
		}
		err = app.render(w, http.StatusOK, "product.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render the Product Page", "template", "product.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		return
	}
	//parse the product quantity and price
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
	//create a new product
	product := &data.Product{
		PName:        productName,
		PQuantity:    productQuantity,
		PPrice:       productPrice,
		PDescription: productDescription,
	}
	//error checker
	err = app.products.Insert(product)
	if err != nil {
		app.logger.Error("failed to insert product into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// VIEW START
// viewHandler is a handler that renders the view page - view.tmpl
func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Current Inventory"
	data.FileInfo = "Here are the products in your inventory."
	err := app.render(w, http.StatusOK, "view.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the view Page", "template", "view.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
