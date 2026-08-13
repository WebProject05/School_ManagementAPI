package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	// Migrating to mux from http
	mux := http.NewServeMux()

	// Note: If you add a trailing slash here (e.g., "/execs/"), Go will automatically
	// send a 301 redirect if the client requests "/execs". This will cause the
	// client to send a second request, hitting our middlewares twice!

	mux.HandleFunc("/", handlers.RootHandlers)

	// GET methods
	mux.HandleFunc("GET /teachers", handlers.TeachersHandler)         // Fetch all or search by query
	mux.HandleFunc("GET /teachers/{id}", handlers.TeachersHandler) // Fetch a specific teacher

	// POST method
	mux.HandleFunc("POST /teachers", handlers.TeachersHandler) // Create new teacher(s)

	// PUT / PATCH methods
	mux.HandleFunc("PUT /teachers/{id}", handlers.TeachersHandler)  // Full replace
	mux.HandleFunc("PATCH /teachers/{id}", handlers.TeachersHandler) // Partial modify

	// DELETE method
	mux.HandleFunc("DELETE /teachers/{id}", handlers.TeachersHandler) // Delete a teacher
	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExcesHandler)

	return mux

}
