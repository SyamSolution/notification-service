package main

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/SyamSolution/notification-service/config"
	"github.com/SyamSolution/notification-service/config/middleware"
	"github.com/SyamSolution/notification-service/internal/consumer"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
)

func main() {
	baseDep := config.NewBaseDep()
	loadEnv(baseDep.Logger)
	db, err := config.NewDbPool(baseDep.Logger)
	if err != nil {
		os.Exit(1)
	}

	dbCollector := middleware.NewStatsCollector("assesment", db)
	prometheus.MustRegister(dbCollector)
	fiberProm := middleware.NewWithRegistry(prometheus.DefaultRegisterer, "transaction-service", "", "", map[string]string{})

	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokers := []string{os.Getenv("KAFKA_BROKER")}
	master, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Panicf("Error creating consumer: %s", err)
	}
	defer func() {
		if err := master.Close(); err != nil {
			log.Panicf("Error closing consumer: %s", err)
		}
	}()

	log.Println("Connected to Kafka broker")

	doneCh := make(chan struct{})
	go consumer.Consumer(master, doneCh)

	app := fiber.New()

	// Define routes and their handlers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//=== metrics route
	fiberProm.RegisterAt(app, "/metrics")
	app.Use(fiberProm.Middleware)

	// Start the Fiber server
	go func() {
		if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
			log.Fatal(err)
		}
	}()

	<-doneCh
}

func loadEnv(logger config.Logger) {
	_, err := os.Stat(".env")
	if err == nil {
		err = godotenv.Load()
		if err != nil {
			logger.Error("no .env files provided")
		}
	}
}
