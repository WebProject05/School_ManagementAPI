package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"restapi/internal/api/middlewares"
	"strings"
	"sync"
)

type user struct {
	Name string `json:"name"`
	Age  string `json:"age"`
	City string `json:"city"`
}

type Teacher struct {
	ID        int
	FirstName string
	LastName  string
	Class     string
	Subject   string
}

var teachers = make(map[int]Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

// Initialize some dummy data
func init() {
	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "10A",
		Subject:   "Physics",
	}
	nextID++

	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "9B",
		Subject:   "Mathematics",
	}
	nextID++

	teachers[nextID] = Teacher{
		ID:        nextID,
		FirstName: "Michael",
		LastName:  "Doe",
		Class:     "11C",
		Subject:   "Chemistry",
	}
	nextID++
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")
	fmt.Println(idStr)

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		teacherList := make([]Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}
	}

	response := struct {
		Status string    `json:"status"`
		Count  int       `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		// A function that handles the get method route
		getTeachersHandler(w, r)

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

	// Note: If you add a trailing slash here (e.g., "/execs/"), Go will automatically
	// send a 301 redirect if the client requests "/execs". This will cause the
	// client to send a second request, hitting our middlewares twice!

	mux.HandleFunc("/", rootHandlers)

	mux.HandleFunc("/teachers/", teachersHandler)

	mux.HandleFunc("/students/", studentsHandler)

	mux.HandleFunc("/execs/", excesHandler)

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// UNComment after development
	// rl := middlewares.NewRateLimiter(5, time.Minute)

	// hppOptions := middlewares.HPPOptions{
	// 	CheckQuery:                  true,
	// 	CheckBody:                   true,
	// 	CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
	// 	WhiteList:                   []string{"sortBy", "sortOrder", "name", "age", "class"},
	// }

	// secureMux := middlewares.Hpp(hppOptions)(rl.Middlware(middlewares.Compression(middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.Cors(mux))))))
	// secureMux := middlewares.ResponseTimeMiddleware(
	// 	middlewares.SecurityHeaders(
	// 		middlewares.Cors(
	// 			rl.Middleware(
	// 				middlewares.Compression(
	// 					middlewares.Hpp(hppOptions)(
	// 						mux,
	// 					),
	// 				),
	// 			),
	// 		),
	// 	),
	// )

	// secureMux := applyMiddlewares(
	// 	mux,
	// 	middlewares.ResponseTimeMiddleware, // 1. Starts timer first
	// 	middlewares.SecurityHeaders,        // 2. Applies headers to all responses
	// 	middlewares.Cors,                   // 3. Handles OPTIONS preflight requests
	// 	rl.Middleware,                      // 4. Blocks spam before heavy processing
	// 	middlewares.Compression,            // 5. Compresses valid, non-blocked payloads
	// 	middlewares.Hpp(hppOptions),        // 6. Sanitizes data right before hitting the app
	// )

	// For development purpose just keep the securit purpose
	secureMux := middlewares.SecurityHeaders(mux)
	// Creating a custom server
	server := &http.Server{
		Addr: port,
		// Handler:   (middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.Cors(mux)))),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server running on the port:", port)
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}

// Middleware is a function that wraps an http.Handler with addional functionality
type Middleware func(http.Handler) http.Handler

func ApplyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middlware := range middlewares {
		handler = middlware(handler)
	}
	return handler
}
