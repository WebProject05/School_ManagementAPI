package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type user struct {
	Name string `json:"name"`
	Age string `json:"age"`
	City string `json:"city"`
}

func main() {
	port := ":3000"

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte("Hello from the server!"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Write([]byte("Read (GET) teachers"))
			return

		case http.MethodPost:
			// Parsing the form data from the user
			err := r.ParseForm()
			if err != nil {
				http.Error(w, "Error Parsing the form", http.StatusBadRequest)
				log.Fatalln("Error parsing the input data:", err)
				return
			}

			// Preparing the response data
			response := make(map[string]interface{})
			for key, value := range r.Form {
				response[key] = value[0]
			}

			// RAW data parsing with json
			body, err := io.ReadAll(r.Body)
			if err != nil {
				return
			}
			defer r.Body.Close()
			fmt.Println("Raw Date:",string(body))


			// UnMarshling the raw body
			var userInstance user
			err = json.Unmarshal(body, &userInstance)
			if err != nil {
				return
			}

			fmt.Println("UnMarshaled JSON:",userInstance)

			fmt.Println("Form data:", r.Form)
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
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the students route!"))
	})

	http.HandleFunc("/execs", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello from the execs route!"))
	})

	fmt.Println("Server running on the port:", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
