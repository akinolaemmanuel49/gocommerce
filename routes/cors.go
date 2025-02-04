package routes

import "github.com/rs/cors"

var corsMiddleware = cors.New(cors.Options{
	AllowedOrigins:   []string{"*"}, // Allow all origins (change in production)
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Authorization", "Content-Type"},
	AllowCredentials: true,
})
