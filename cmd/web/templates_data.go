package main

type TemplateData struct {
	Title      string
	HeaderText string
	FormErrors map[string]string
	FormData   map[string]string
	// Feedbacks  []*data.Feedback // Feedbacks is a slice of pointers to data.Feedback
	// Journals   []*data.Journal  //Journal
	// Todos      []*data.Todo     //Todo
}

func NewTemplateData() *TemplateData {
	return &TemplateData{
		Title:      "Default Title",
		HeaderText: "Default HeaderText",

		FormErrors: map[string]string{},
		FormData:   map[string]string{},
		// Feedbacks:  []*data.Feedback{}, // Initialize the slice
		// Journals:   []*data.Journal{},  //Journal
		// Todos:      []*data.Todo{},     //Initialize the todo slice

	}
}
