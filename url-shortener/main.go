package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	h "github.com/NoCapCbas/url-shortener/api"
	mongoRepo "github.com/NoCapCbas/url-shortener/repository/mongo"
	redisRepo "github.com/NoCapCbas/url-shortener/repository/redis"
	inMemRepo "github.com/NoCapCbas/url-shortener/repository/inMemory"

)

// https://www.google.com -> 98sj1-293
// http://localhost:8000/98sj1-293 -> https://www.google.com

// repo <- service -> serializer  -> http

func main() {
	repo := chooseRepo()
	service := urlshortener.NewRedirectService(repo)
	handler := h.NewHandler(service)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(httpPort(), r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)

}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo() urlshortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
    log.Println("Using redis for data source...")
		redisURL := os.Getenv("REDIS_URL")
		repo, err := redisRepo.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
    log.Println("Using mongo for data source...")
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongoRepo.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
  default:
    log.Println("Using in memory for data source...")
    repo := inMemRepo.NewInMemoryRepository() 
	}
	return nil
}
