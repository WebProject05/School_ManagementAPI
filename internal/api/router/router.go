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

	mux.HandleFunc("GET /teachers", handlers.GetTeachersHandler)      // Fetch all or search by query
	mux.HandleFunc("POST /teachers", handlers.PostTeacherHandler) // Create new teacher(s)
	mux.HandleFunc("PATCH /teachers/", handlers.PatchTeacherHandler) // Partial modify
	mux.HandleFunc("DELETE /teachers/", handlers.DeleteTeacherHandler) // Delete a teacher
	
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler) // Fetch a specific teacher
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)  // Full replace
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler) // Partial modify
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler) // Delete a teacher

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExcesHandler) // <-- Add the function name back
	return mux

}
