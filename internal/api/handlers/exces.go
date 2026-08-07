package handlers

import "net/http"

func ExcesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write([]byte("Read (GET) exces"))
		return

	case http.MethodPost:
		w.Write([]byte("Create (POST) exces"))
		return

	case http.MethodPut:
		w.Write([]byte("Update (PUT) exces"))
		return

	case http.MethodPatch:
		w.Write([]byte("Partial Update (PATCH) exces"))
		return

	case http.MethodDelete:
		w.Write([]byte("Delete (DELETE) exces"))
		return

	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
	w.Write([]byte("Hello from the execs route!"))
}
