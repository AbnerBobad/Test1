package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/AbnerBobad/final_project/internal/data"
	"github.com/AbnerBobad/final_project/internal/validator"
	"github.com/justinas/nosurf"
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
	data.CSRFToken = nosurf.Token(r)
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
	data.CSRFToken = nosurf.Token(r)

	err := app.render(w, http.StatusOK, "product.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render the Product Page", "template", "product.tmpl", "error", err, "url", r.URL.Path, "method", r.Method)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// product
func (app *application) createProduct(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err := r.ParseForm()
	if err != nil {
		app.logger.Error("failed to parse products form data", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	productName := r.PostForm.Get("product_name")
	productQuantityStr := r.PostForm.Get("product_quantity")
	productPriceStr := r.PostForm.Get("product_price")
	productDescription := r.PostForm.Get("product_description")

	productQuantity, err := strconv.ParseInt(productQuantityStr, 10, 64)

	productPrice, err := strconv.ParseFloat(productPriceStr, 64)

	userID, ok := app.session.Get(r, "authenticateUserID").(int64)
	if !ok {
		app.logger.Error("failed to retrieve user ID from session")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	product := &data.Product{
		PName:        productName,
		PQuantity:    productQuantity,
		PPrice:       productPrice,
		PDescription: productDescription,
		User:         userID,
	}

	v := validator.NewValidator()
	data.ValidateProduct(v, product)
	if !v.ValidData() {
		formData := NewTemplateData()
		formData.Title = "StockTrack"
		formData.HeaderText = "Add New Products"
		formData.FileInfo = "Please fill in the product details below."
		formData.FormErrors = v.Errors
		formData.IsAuthenticated = app.IsAuthenticated(r)
		formData.CSRFToken = nosurf.Token(r)
		formData.FormData = map[string]string{
			"product_name":        productName,
			"product_quantity":    productQuantityStr,
			"product_price":       productPriceStr,
			"product_description": productDescription,
		}
		err = app.render(w, http.StatusOK, "product.tmpl", formData)
		if err != nil {
			app.logger.Error("failed to render the Product Page", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	err = app.products.Insert(product)
	if err != nil {
		app.logger.Error("failed to insert product into database", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/product?submitted=true", http.StatusSeeOther)
}

// view
func (app *application) listUserProducts(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	userID, ok := app.session.Get(r, "authenticateUserID").(int64)
	if !ok {
		app.logger.Error("failed to retrieve user ID from session")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	products, err := app.products.GetAllForUser(userID)
	if err != nil {
		app.logger.Error("failed to retrieve products for user", "error", err)
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
	data.Submitted = submitted
	data.IsAuthenticated = true
	data.Products = products
	data.CSRFToken = nosurf.Token(r)

	err = app.render(w, http.StatusOK, "view.tmpl", data)
	if err != nil {
		app.logger.Error("failed to render product list", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// func (app *application) editProductForm(w http.ResponseWriter, r *http.Request) {
// 	// Check if user is authenticated
// 	if !app.IsAuthenticated(r) {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Get the product ID from the URL query parameter
// 	idStr := r.URL.Query().Get("id")
// 	id, err := strconv.ParseInt(idStr, 10, 64)
// 	if err != nil || id <= 0 {
// 		http.NotFound(w, r)
// 		return
// 	}

// 	// Retrieve the product by ID
// 	product, err := app.products.GetByID(id)
// 	if err != nil {
// 		app.logger.Error("error getting product", "error", err)
// 		http.Error(w, "Product not found", http.StatusNotFound)
// 		return
// 	}

// 	// Retrieve the logged-in user's ID from the session
// 	userID, ok := app.session.Get(r, "authenticateUserID").(int64)
// 	if !ok {
// 		app.logger.Error("failed to retrieve userID from session", "error", err)
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Ensure the logged-in user is the owner of the product
// 	if product.User != userID {
// 		app.logger.Error("user does not own the product", "productUser", product.User, "userID", userID)
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}

// 	// Prepare data to render the edit form
// 	data := NewTemplateData()
// 	data.Title = "Edit Product"
// 	data.HeaderText = "Edit Product"
// 	data.Product = product
// 	data.IsAuthenticated = app.IsAuthenticated(r)
// 	data.CSRFToken = nosurf.Token(r)

// 	// Render the edit form template
// 	err = app.render(w, http.StatusOK, "edit_product.tmpl", data)
// 	if err != nil {
// 		app.logger.Error("failed to render edit form", "error", err)
// 		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 	}
// }

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
	data.CSRFToken = nosurf.Token(r)

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
	// http.Redirect(w, r, "/view?submitted=true", http.StatusSeeOther)
	// http.Redirect(w, r, fmt.Sprintf("/edit?id=%d&submitted=true", id), http.StatusSeeOther)
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

// // search
func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {
	if !app.IsAuthenticated(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the current user ID from the session
	userID, ok := app.session.Get(r, "authenticateUserID").(int64)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the search query
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Redirect(w, r, "/view", http.StatusSeeOther)
		return
	}

	// Perform a filtered search for products owned by the user and matching the query
	products, err := app.products.SearchByUser(query, userID)
	if err != nil {
		app.logger.Error("search failed", "error", err)
		http.Error(w, "Search error", http.StatusInternalServerError)
		return
	}

	// Add stock warnings
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
	data.CSRFToken = nosurf.Token(r)

	err = app.render(w, http.StatusOK, "view.tmpl", data)
	if err != nil {
		app.logger.Error("render search results failed", "error", err)
		http.Error(w, "Render error", http.StatusInternalServerError)
	}
}

// func (app *application) searchProducts(w http.ResponseWriter, r *http.Request) {
// 	if !app.IsAuthenticated(r) {
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		return
// 	}
// 	query := r.URL.Query().Get("query")
// 	if query == "" {
// 		http.Redirect(w, r, "/view", http.StatusSeeOther)
// 		return
// 	}

// 	products, err := app.products.Search(query)
// 	if err != nil {
// 		app.logger.Error("search failed", "error", err)
// 		http.Error(w, "Search error", http.StatusInternalServerError)
// 		return
// 	}
// 	//warning
// 	for _, product := range products {
// 		switch {
// 		case product.PQuantity == 0:
// 			product.StockStatus = "Out of Stock"
// 		case product.PQuantity <= 5:
// 			product.StockStatus = "Stock Low"
// 		default:
// 			product.StockStatus = "Available"
// 		}
// 	}

// 	data := NewTemplateData()
// 	data.Products = products
// 	data.Title = "Search Results"
// 	data.HeaderText = "Search Results for: " + query
// 	data.IsAuthenticated = app.IsAuthenticated(r)
// 	data.CSRFToken = nosurf.Token(r)

// 	err = app.render(w, http.StatusOK, "view.tmpl", data)
// 	if err != nil {
// 		app.logger.Error("render search results failed", "error", err)
// 		http.Error(w, "Render error", http.StatusInternalServerError)
// 	}
// }

// Login page handler
func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	data := NewTemplateData()
	data.Title = "Login"
	data.HeaderText = "Login Page"
	data.FileInfo = "Please login to continue."
	data.IsAuthenticated = app.IsAuthenticated(r)
	data.CSRFToken = nosurf.Token(r)
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
	data.IsAuthenticated = app.IsAuthenticated(r)
	data.CSRFToken = nosurf.Token(r)

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
		data.IsAuthenticated = app.IsAuthenticated(r)
		data.CSRFToken = nosurf.Token(r)
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
	data.IsAuthenticated = app.IsAuthenticated(r)
	data.CSRFToken = nosurf.Token(r)
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
		data.IsAuthenticated = app.IsAuthenticated(r)
		data.CSRFToken = nosurf.Token(r)
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
			data.IsAuthenticated = app.IsAuthenticated(r)
			data.CSRFToken = nosurf.Token(r)
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
