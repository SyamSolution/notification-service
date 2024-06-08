package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/IBM/sarama"
	"github.com/SyamSolution/notification-service/config"
	"github.com/SyamSolution/notification-service/config/middleware"
	"github.com/SyamSolution/notification-service/internal/consumer"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
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

	go consumer.StartConsumer()
	app := fiber.New()

	app.Use(recover.New())
	app.Use(cors.New())
	app.Use(pprof.New())
	app.Use(logger.New(logger.Config{
		// Format: "[${time}] ${status} - ${latency} ${method} ${path}\n",
		Format:       `${time} {"router_activity" : [${status},"${latency}","${method}","${path}"], "query_param":${queryParams}, "body_param":${body}}` + "\n",
		TimeInterval: time.Millisecond,
		TimeFormat:   "02-01-2006 15:04:05",
		TimeZone:     "Indonesia/Jakarta",
	}))

	// Define routes and their handlers
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	//=== metrics route
	fiberProm.RegisterAt(app, "/metrics")
	app.Use(fiberProm.Middleware)

	// Start the Fiber server
	if err := app.Listen(fmt.Sprintf(":%s", os.Getenv("APP_PORT"))); err != nil {
		log.Fatal(err)
	}
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
