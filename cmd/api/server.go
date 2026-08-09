package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

var teachers = make(map[int]models.Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

func main() {
	// Load environment variables from .env file
	_ = godotenv.Load()

	db, err := sqlconnect.ConnectDb()
	if err != nil {
		fmt.Println("Error (loc: server.go)", err)
		return
	}

	fmt.Printf("Database: %v \n", db)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = ":3000"
	}

	cert := "cert.pem"
	key := "key.pem"

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	// UNComment after development
	rl := middlewares.NewRateLimiter(50, time.Minute)

	// Update this HPP Options often when made changes to the handlers and middlewares
	hppOptions := middlewares.HPPOptions{
		CheckQuery:                  true,
		CheckBody:                   true,
		CheckBodyOnlyForContentType: "application/x-www-form-urlencoded",
		WhiteList:                   []string{"sortBy", "sortOrder", "name", "age", "class", "email"},
	}

	// secureMux := middlewares.Hpp(hppOptions)(rl.Middlware(middlewares.Compression(middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.Cors(mux))))))
	// secureMux := middlewares.ResponseTimeMiddleware(
	//  middlewares.SecurityHeaders(
	//      middlewares.Cors(
	//          rl.Middleware(
	//              middlewares.Compression(
	//                  middlewares.Hpp(hppOptions)(
	//                      mux,
	//                  ),
	//              ),
	//          ),
	//      ),
	//  ),
	// )
	router := router.Router()

	secureMux := utils.ApplyMiddlewares(
		router,
		middlewares.ResponseTimeMiddleware, // 1. Starts timer first
		middlewares.SecurityHeaders,        // 2. Applies headers to all responses
		middlewares.Cors,                   // 3. Handles OPTIONS preflight requests
		rl.Middleware,                      // 4. Blocks spam before heavy processing
		middlewares.Compression,            // 5. Compresses valid, non-blocked payloads
		middlewares.Hpp(hppOptions),        // 6. Sanitizes data right before hitting the app
	)

	// For development purpose just keep the securit purpose
	// secureMux := middlewares.SecurityHeaders(router)
	// Creating a custom server
	server := &http.Server{
		Addr: port,
		// Handler:   (middlewares.ResponseTimeMiddleware(middlewares.SecurityHeaders(middlewares.Cors(mux)))),
		Handler:   secureMux,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server running on the port:", port)
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Fatalln("Error starting the server:", err)
	}
}
