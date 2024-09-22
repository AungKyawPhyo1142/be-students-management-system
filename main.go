package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/controllers"
	"github.com/AungKyawPhyo1142/be-students-management-system/handlers"
	"github.com/AungKyawPhyo1142/be-students-management-system/migrations"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	port := os.Getenv("PORT")

	// connect to database
	config.ConnectDB()

	// run migrations
	migrations.Migrate()

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
	v1Router.Post("/user", controllers.CreateUser)
	v1Router.Get("/user", controllers.GetAllUsers)

	// Student Related Routes
	v1Router.Post("/students", controllers.CreateStudent)        // create
	v1Router.Patch("/students/{id}", controllers.EditStudent)    // update
	v1Router.Delete("/students/{id}", controllers.DeleteStudent) // delete student by id
	v1Router.Get("/students/{id}", controllers.GetStudentByID)   // get student by id
	v1Router.Get("/students", controllers.GetAllStudents)        // get all students

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
