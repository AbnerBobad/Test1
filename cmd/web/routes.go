package main

import (
	"net/http"
)

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	//static files
	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("GET /static/", http.StripPrefix("/static", fileServer))

	// //Routes Definitions (GET)
	// mux.HandleFunc("GET /", app.home)                        //display home page
	// mux.HandleFunc("GET /feedback/{$}", app.feedbackHandler) //to submit feedback
	// mux.HandleFunc("GET /journal/{$}", app.journalHandler)   // to submit form
	// mux.HandleFunc("GET /todo/{$}", app.todoHandler)         //to submit todo

	// //Success Pages
	// mux.HandleFunc("GET /feedback/success", app.feedbackSuccess) //feedback success
	// mux.HandleFunc("GET /journal/success", app.journalSuccess)   //Journal Success Route
	// mux.HandleFunc("GET /todo/success", app.todoSuccess)         //Todo Success Route

	// //Retrieve Data GETS
	// mux.HandleFunc("GET /feedbacks", app.viewFeedbacks) //view feedbacks
	// mux.HandleFunc("GET /journals", app.viewJournals)   //view journals
	// mux.HandleFunc("GET /todos", app.viewTodos)         //view todos

	// //Routes Definitions (POST)
	// mux.HandleFunc("POST /feedback/success", app.createFeedback) //handle feedback submission
	// mux.HandleFunc("POST /journal/success", app.createJournal)   //handle journal submission
	// mux.HandleFunc("POST /todo/success", app.createTodo)         //handle todo submission

	return app.loggingMiddleware(mux)
}
