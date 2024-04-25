package main

import (
	"log"
	"net/http"

	"github.com.br/gibranct/golang-course-api/configs"
	"github.com.br/gibranct/golang-course-api/internal/entity"
	"github.com.br/gibranct/golang-course-api/internal/infra/database"
	"github.com.br/gibranct/golang-course-api/internal/infra/webservers/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, *configs.TokenAuth, configs.JWTExpiresIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(LogRequest)

	r.Route("/products", func(inRoute chi.Router) {
		inRoute.Use(jwtauth.Verifier(configs.TokenAuth))
		inRoute.Use(jwtauth.Authenticator)
		inRoute.Post("/", productHandler.CreateProduct)
		inRoute.Get("/{id}", productHandler.FindByID)
		inRoute.Put("/{id}", productHandler.Update)
		inRoute.Delete("/{id}", productHandler.DeleteById)
	})

	r.Post("/login", userHandler.Login)
	r.Post("/users", userHandler.CreateUser)

	http.ListenAndServe(":8000", r)
}

func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
