package main

import (
	"github.com/AbnerBobad/final_project/internal/data"
)

type TemplateData struct {
	Title      string
	HeaderText string
	FileInfo   string

	FormErrors map[string]string
	FormData   map[string]string

	Products []*data.Product // Products is a slice of pointers to data.Product
	Users    []*data.User

	Submitted bool
	Product   *data.Product
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FileInfo:   "Default FileInfo",

		FormErrors: map[string]string{},
		FormData:   map[string]string{},

		// Initialize the slice
		Products: []*data.Product{},
		Users:    []*data.User{},

		Submitted: false,
		Product:   nil,
	}
}
