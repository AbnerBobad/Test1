package main

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// Test handler for Login page
// func TestLoginHandler(t *testing.T) {
// 	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
// 	app := &application{
// 		logger: logger,
// 	}
// 	req := httptest.NewRequest("GET", "/", nil)
// 	rr := httptest.NewRecorder()
// 	handler := http.HandlerFunc(app.loginHandler)
// 	handler.ServeHTTP(rr, req)
// 	status := rr.Code
// 	if status != http.StatusOK {
// 		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
// 	}
// 	expected := "Login Page\n"
// 	got := rr.Body.String()
// 	if got != expected {
// 		t.Errorf("got %q, expected %q", got, expected)
// 	}
// }

// Test handler for main page
func TestHomeHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.mainHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}
	expected := "Welcome to StockTrack\n"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// Test handler for the product page
func TestProductHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	req := httptest.NewRequest("GET", "/product", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.mainHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}
	expected := "Add New Products\n"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

// Test handler for the view page
func TestViewHandler(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	req := httptest.NewRequest("GET", "/view", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.mainHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}
	expected := "Current Inventory\n"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}

func TestEditProduct(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &application{
		logger: logger,
	}
	req := httptest.NewRequest("GET", "/product/edit", nil)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(app.mainHandler)
	handler.ServeHTTP(rr, req)
	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("got status code %v, expected status code %v.", status, http.StatusOK)
	}
	expected := "Edit Product\n"
	got := rr.Body.String()
	if got != expected {
		t.Errorf("got %q, expected %q", got, expected)
	}
}
