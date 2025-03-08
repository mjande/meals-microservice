package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
	_ "github.com/joho/godotenv/autoload"
	"github.com/mjande/meals-microservice/database"
	"github.com/mjande/meals-microservice/handlers"
)

func main() {
	// Connect to database
	err := database.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer database.DB.Close()

	// Create new router
	router := chi.NewRouter()

	// CORS Middleware
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{os.Getenv("CLIENT_URL")},
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}))

	// Logging Middleware
	router.Use(middleware.Logger)

	// JWT Middleware
	tokenAuth := jwtauth.New("HS256", []byte(os.Getenv("SECRET_KEY")), nil)

	router.Use(jwtauth.Verifier(tokenAuth))

	// Routes
	router.Route("/meals", func(r chi.Router) {
		r.Get("/", handlers.GetMeals)
		r.Post("/", handlers.PostMeal)
		r.Delete("/{id}", handlers.DeleteMeal)
	})

	// Start server
	log.Printf("Meal service listening on port %s", os.Getenv("PORT"))
	err = http.ListenAndServe(":"+os.Getenv("PORT"), router)
	if err != nil {
		log.Fatal(err)
	}
}
