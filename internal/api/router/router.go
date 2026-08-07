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

	mux.HandleFunc("/teachers/", handlers.TeachersHandler)

	mux.HandleFunc("/students/", handlers.StudentsHandler)

	mux.HandleFunc("/execs/", handlers.ExcesHandler)

	return mux

}
