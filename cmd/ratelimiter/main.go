package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/leonardosm2/Rate-Limiter/configs"
	"github.com/leonardosm2/Rate-Limiter/internal/infra/web"
	"github.com/leonardosm2/Rate-Limiter/internal/middleware"
	"github.com/leonardosm2/Rate-Limiter/internal/repository"
	"github.com/leonardosm2/Rate-Limiter/internal/usecase"
)

func main() {
	configs, err := configs.LoadConfig("../..")
	if err != nil {
		panic(err)
	}

	rateLimiterRepository := repository.NewRedisRepository(
		configs.RedisHost,
		configs.RedisPort,
		configs.RedisDb,
	)
	rateLimiterUseCase := usecase.NewRateLimiterUseCase(
		rateLimiterRepository,
		configs.RateLimitDefault,
		configs.TimeBlockDefault,
	)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(
		*rateLimiterUseCase,
	)

	server := web.NewServer(*rateLimiterMiddleware)
	router := server.CreateServer()

	fmt.Println("Starting web server on port ", configs.WebServerPort)
	log.Fatal(http.ListenAndServe(configs.WebServerPort, router))
}
