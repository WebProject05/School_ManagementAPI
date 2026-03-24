package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

func rootHandlers(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Write([]byte("Hello from the server!"))
}

func teachersHandler(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			fmt.Println(r.URL.Path)
			path := strings.TrimPrefix(r.URL.Path, "/teachers/")
			userID := strings.TrimSuffix(path, "/")

			fmt.Println("User ID:", userID)


			fmt.Println("Query Params:", r.URL.Query())
			queryParam := r.URL.Query()
			name := queryParam.Get("name")
			age := queryParam.Get("age")
			fmt.Println("Name from the query:", name)
			fmt.Println("Age from the Query:", age)
			w.Write([]byte("Read (GET) teachers"))
			return

		case http.MethodPost:
			w.Write([]byte("Create (POST) teacher"))
			return

		case http.MethodPut:
			w.Write([]byte("Update (PUT) teacher"))
			return

		case http.MethodPatch:
			w.Write([]byte("Partial Update (PATCH) teacher"))
			return

		case http.MethodDelete:
			w.Write([]byte("Delete (DELETE) teacher"))
			return

		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
		w.Write([]byte("Hello from the teachers route!"))
	}


func studentsHandler(w http.ResponseWriter, r *http.Request) {
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



func excesHandler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"

	// Migrating to mux from http
	mux := http.NewServeMux()

	mux.HandleFunc("/", rootHandlers)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", excesHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// Creating a custom server
	server := &http.Server{
		Addr: port,
		Handler: mux ,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server running on the port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
