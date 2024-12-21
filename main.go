package main

import (
	"github.com/AungKyawPhyo1142/be-students-management-system/helpers"
	"log"
	"net/http"
	"os"

	"github.com/AungKyawPhyo1142/be-students-management-system/config"
	"github.com/AungKyawPhyo1142/be-students-management-system/controllers"
	"github.com/AungKyawPhyo1142/be-students-management-system/handlers"
	"github.com/AungKyawPhyo1142/be-students-management-system/middleware"
	"github.com/AungKyawPhyo1142/be-students-management-system/migrations"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
)

func main() {

	godotenv.Load(".env")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

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
	v1Router.Get("/user", controllers.GetAllUsers)

	// auth
	v1Router.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", controllers.Register) // register
		auth.Post("/login", controllers.Login)       // login
	})

	// routes protected with AuthMiddleware
	v1Router.Route("/", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)

		// student
		r.Route("/student", func(r chi.Router) {
			r.Post("/", controllers.CreateStudent)       // create
			r.Patch("/{id}", controllers.EditStudent)    // update
			r.Delete("/{id}", controllers.DeleteStudent) // delete student by id
			r.Get("/{id}", controllers.GetStudentByID)   // get student by id
			r.Get("/", controllers.GetAllStudents)       // get all students

			r.Post("/{id}/{code}", controllers.AssignStudentToClass)     // assign student to class
			r.Get("/{code}", controllers.GetAllStudentsFromClass)        // get students from class
			r.Delete("/{id}/{code}", controllers.DeleteStudentFromClass) // remove student from class
		})

		// class
		r.Route("/class", func(r chi.Router) {
			r.Post("/", controllers.CreateClass)         // create
			r.Get("/", controllers.GetAllClasses)        // get all classes
			r.Get("/{code}", controllers.GetClassByID)   // get class by class_code
			r.Patch("/{code}", controllers.UpdateClass)  // update
			r.Delete("/{code}", controllers.DeleteClass) // delete class by class_code
		})

		// instructor
		r.Route("/instructor", func(r chi.Router) {
			r.Post("/", controllers.CreateInstructor)       // create
			r.Get("/", controllers.GetAllInstructors)       // get all instructors
			r.Get("/{id}", controllers.GetInstructorByID)   // get instructor by id
			r.Patch("/{id}", controllers.EditInstructor)    // update instructor
			r.Delete("/{id}", controllers.DeleteInstructor) // delete instructor by id

			r.Post("/{id}/{code}", controllers.AssignInstructorToClass)     // assign instructor to a class
			r.Get("/{code}", controllers.GetAllInstructorsFromClass)        // get all instructors from class
			r.Delete("/{id}/{code}", controllers.DeleteInstructorFromClass) // remove instructor from class
		})

	})

	// mount the v1 router to main/default router
	router.Mount("/v1", v1Router)

	// print all available routes
	helpers.PrintRoutes(router)

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
