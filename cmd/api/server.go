package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"restapi/internal/api/middlewares"
	"restapi/internal/api/router"
	"restapi/internal/models"
	"sync"
)

var teachers = make(map[int]models.Teacher)
var mutex = &sync.Mutex{}
var nextID = 1

func main() {
	port := ":3000"

	cert := "cert.pem"
	key := "key.pem"



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

	// secureMux := utils.ApplyMiddlewares(
	// 	mux,
	// 	middlewares.ResponseTimeMiddleware, // 1. Starts timer first
	// 	middlewares.SecurityHeaders,        // 2. Applies headers to all responses
	// 	middlewares.Cors,                   // 3. Handles OPTIONS preflight requests
	// 	rl.Middleware,                      // 4. Blocks spam before heavy processing
	// 	middlewares.Compression,            // 5. Compresses valid, non-blocked payloads
	// 	middlewares.Hpp(hppOptions),        // 6. Sanitizes data right before hitting the app
	// )

	// For development purpose just keep the securit purpose
	router := router.Router()
	secureMux := middlewares.SecurityHeaders(router)
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
