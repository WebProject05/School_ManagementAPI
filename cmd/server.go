package main

import (
	"fmt"
	"log"
	"net/http"
)

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
		if r.Method == http.MethodGet {
			w.Write([]byte("Hello GET Method on the teachers route"))
			return
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