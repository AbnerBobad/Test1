package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AbnerBobad/final_project/internal/data"
	"github.com/AbnerBobad/final_project/internal/validator"
	"golang.org/x/crypto/bcrypt"
)

// MAIN PAGE START
// mainHandler is a handler that renders the main page - main.tmpl
func (app *application) mainHandler(w http.ResponseWriter, r *http.Request) {
	submitted := r.URL.Query().Get("submitted") == "true"
	data := NewTemplateData()
	data.Title = "StockTrack"
	data.HeaderText = "Welcome to StockTrack"
	data.FileInfo = "Manage your inventory efficiently and stay updated on stock levels."
	data.Submitted = submitted
	data.IsAuthenticated = app.IsAuthenticated(r)
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
	submitted := r.URL.Query().Get("submitted") == "true"

	data.Title = "StockTrack"
	data.HeaderText = "Add New Products"
	data.FileInfo = "Please fill in the product details below."
	data.Submitted = submitted
	data.IsAuthenticated = app.IsAuthenticated(r)

	err := app.render(w, http.StatusOK, "product.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Product Page", "template", "product.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Product form creation2
func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	//testing guard
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//parsed data form
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse products form data", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	//parse data member
	productName := r.PostForm.Get("product_name")
	productQuantityStr := r.PostForm.Get("product_quantity")
	productPriceStr := r.PostForm.Get("product_price")
	productDescription := r.PostForm.Get("product_description")

	//converted data members
	productQuantity, err := strconv.ParseInt(productQuantityStr, 10, 64)
	productPrice, err := strconv.ParseFloat(productPriceStr, 64)

	product := &data.Product{
		PName:        productName,
		PQuantity:    productQuantity,
		PPrice:       productPrice,
		PDescription: productDescription,
	}
	//validate data
	v := validator.NewValidator()
	data.ValidateProduct(v, product)
	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "StockTrack"
		data.HeaderText = "Add New Products"
		data.FileInfo = "Please fill in the product details below."
		data.FormErrors = v.Errors
		data.IsAuthenticated = app.IsAuthenticated(r)
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
	//error checker
	err = app.products.Insert(product)
	if err != nil {
		app.logger.Error("failed to insert product into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	//set a session data

	http.Redirect(w, r, "/product?submitted=true", http.StatusSeeOther)

}

// VIEW START
// viewHandler is a handler that renders the view page - view.tmpl
func (app *application) viewHandler(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	//get all products from the database
	products, err := app.products.GetAll()
	if err != nil {
		app.logger.Error("failed to get products from database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	//warning
	for _, product := range products {
		switch {
		case product.PQuantity == 0:
			product.StockStatus = "Out of Stock"
		case product.PQuantity <= 5:
			product.StockStatus = "Stock Low"
		default:
			product.StockStatus = "Available"
		}
	}

	data := NewTemplateData()
	submitted := r.URL.Query().Get("submitted") == "true"
	data.Title = "StockTrack"
	data.HeaderText = "Current Inventory"
	data.FileInfo = "Here are the products in your inventory."
	data.Products = products
	data.Submitted = submitted
	data.IsAuthenticated = app.IsAuthenticated(r)

	err = app.render(w, http.StatusOK, "view.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the view Page", "template", "view.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Edit
func (app *application) editProductForm(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	product, err := app.products.GetByID(id)
	if err != nil {
		app.logger.Error("error getting product", "error", err)
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	data := NewTemplateData()

	data.Title = "Edit Product"
	data.HeaderText = "Edit Product"
	data.Product = product
	data.IsAuthenticated = app.IsAuthenticated(r)

	err = app.render(w, http.StatusOK, "edit_product.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render edit form", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// update
func (app *application) updateProduct(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Invalid form", http.StatusBadRequest)
		return
	}

	idStr := r.PostFormValue("product_id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	product := &data.Product{
		ID:           id,
		PName:        r.PostFormValue("product_name"),
		PDescription: r.PostFormValue("product_description"),
	}

	product.PQuantity, err = strconv.ParseInt(r.PostFormValue("product_quantity"), 10, 64)
	product.PPrice, err = strconv.ParseFloat(r.PostFormValue("product_price"), 64)

	err = app.products.Update(product)
	if err != nil {
		app.logger.Error("failed to update product", "error", err)
		http.Error(w, "Could not update product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view?submitted=true", http.StatusSeeOther)
	http.Redirect(w, r, fmt.Sprintf("/edit?id=%d&submitted=true", id), http.StatusSeeOther)
}

// delete
func (app *application) deleteProduct(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	idStr := r.URL.Query().Get("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil || id <= 0 {
		http.NotFound(w, r)
		return
	}

	err = app.products.Delete(id)
	if err != nil {
		app.logger.Error("failed to delete product", "error", err)
		http.Error(w, "Could not delete product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view", http.StatusSeeOther)
}

// search
func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}

	products, err := app.products.Search(query)
	if err != nil {
		app.logger.Error("search failed", "error", err)
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}
	//warning
	for _, product := range products {
		switch {
		case product.PQuantity == 0:
			product.StockStatus = "Out of Stock"
		case product.PQuantity <= 5:
			product.StockStatus = "Stock Low"
		default:
			product.StockStatus = "Available"
		}
	}

	data := NewTemplateData()
	data.Products = products
	data.Title = "Search Results"
	data.HeaderText = "Search Results for: " + query
	data.IsAuthenticated = app.IsAuthenticated(r)

	err = app.render(w, http.StatusOK, "view.tmpl", data)
	if err != nil {
		app.logger.Error("render search results failed", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

// Login page handler
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Login"
	data.HeaderText = "Login Page"
	data.FileInfo = "Please login to continue."
	err := app.render(w, http.StatusOK, "login.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Login page", "template", "login.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse users form data", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password_hash") // Use "password" as the field name

	// Validate input
	errors_user := make(map[string]string)

	data := NewTemplateData()
	data.Title = "Login"
	data.HeaderText = "Login Page"
	data.FileInfo = "Please login to continue."
	data.FormErrors = errors_user

	data.FormData = map[string]string{
		"email": email,
	}

	// Authenticate
	id, err := app.users.Authenticate(email, password)
	if err != nil {
		data := NewTemplateData()
		data.Title = "Login"
		data.HeaderText = "Login Page"
		data.FileInfo = "Please login to continue."
		data.FormErrors = map[string]string{"default": "Email or password is incorrect"}
		data.FormData = map[string]string{
			"email": email,
		}

		if rErr := app.render(w, http.StatusOK, "login.tmpl", data); rErr != nil {
			app.logger.Error("failed to render login after auth error", "err", rErr)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}
	app.session.Put(r, "authenticateUserID", id)

	// Redirect to products page on successful login
	http.Redirect(w, r, "/product", http.StatusSeeOther)
}

// Sign Up
func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Sign Up"
	data.HeaderText = "Create a new account"
	data.FileInfo = "Please fill in the form to create a new account."
	err := app.render(w, http.StatusOK, "signup.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Signup page", "template", "signup.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	// Parse the form data
	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse users form data", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Get form values
	username := r.PostForm.Get("username")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password_hash")

	// Validate data
	user := &data.User{
		UName:  username,
		Email:  email,
		URole:  "admin",
		Active: true,
	}

	v := validator.NewValidator()
	data.ValidateUser(v, &data.User{
		UName:          username,
		Email:          email,
		HashedPassword: []byte(password),
	})

	if !v.ValidData() {
		data := NewTemplateData()
		data.Title = "Sign Up"
		data.HeaderText = "Create a new account"
		data.FileInfo = "Please fill in the form to create a new account."
		data.FormErrors = v.Errors
		data.FormData = map[string]string{
			"username": username,
			"email":    email,
		}

		err := app.render(w, http.StatusOK, "signup.tmpl", data)
		if err != nil {
			app.logger.Error("failed to render the signup page", "template", "signup.tmpl", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.logger.Error("failed to hash password", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.HashedPassword = hashedPassword

	if err := app.users.Insert(user); err != nil {
		if errors.Is(err, data.ErrDuplicateEmail) {
			// Re-render with friendly message
			data := NewTemplateData()
			data.Title = "Sign Up"
			data.FormErrors = map[string]string{"email": "Email is already registered"}
			data.FormData = map[string]string{
				"username": username,
				"email":    email,
			}

			if rErr := app.render(w, http.StatusOK, "signup.tmpl", data); rErr != nil {
				app.logger.Error("render signup tmpl", "err", rErr)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
			return
		}
		app.logger.Error("Insert User error", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// Logout
func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	//clear the session
	app.session.Remove(r, "authenticateUserID")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
