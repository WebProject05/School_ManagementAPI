package handlers

import "net/http"

func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Read (GET) students"))
		return

	case http.MethodPost:
		w.Write([]byte("Create (POST) students"))
		return

	case http.MethodPut:
		w.Write([]byte("Update (PUT) students"))
		return

	case http.MethodPatch:
		w.Write([]byte("Partial Update (PATCH) students"))
		return

	case http.MethodDelete:
		w.Write([]byte("Delete (DELETE) students"))
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.Write([]byte("Hello from the students route!"))
}
