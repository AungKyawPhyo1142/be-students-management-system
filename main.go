package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AungKyawPhyo1142/be-students-management-system/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	// define a default router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		MaxAge:           300,
		AllowCredentials: false,
	}))

	// define v1 Router
	v1Router := chi.NewRouter()
	v1Router.Get("/ready", handlers.HandlerReady)

	// mount the v1 router to main/default router
	router.Mount("/v1", v1Router)

	// define server with apiRoutes and port number to listen
	server := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server is listening on port %v", port)
	ServerErr := server.ListenAndServe()

	if ServerErr != nil {
		log.Fatal(ServerErr)
	}

}
