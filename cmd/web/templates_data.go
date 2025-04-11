package main

type TemplateData struct {
	Title      string
	HeaderText string
	FileInfo   string
	FormErrors map[string]string
	FormData   map[string]string
	Products   []*data.Product // Products is a slice of pointers to data.Product
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",
		FileInfo:   "Default FileInfo",
		FormErrors: map[string]string{},
		FormData:   map[string]string{},
		Products:   []*data.Product{}, // Initialize the slice

	}
}
