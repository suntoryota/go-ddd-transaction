package main

import (
	"fmt"
	"log"

	"onlineShop/auth"
	"onlineShop/config"
	"onlineShop/database"
	"onlineShop/middleware"
	"onlineShop/product"
	"onlineShop/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	filename := "cmd/api/config.yaml"
	if err := config.LoadConfig(filename); err != nil {
		panic(err)
	}

	db, err := database.ConnectPostgres(config.Cfg.DB)
	if err != nil {
		panic(err)
	}

	if db != nil {
		log.Println("db connected")
	}

	router := fiber.New(fiber.Config{
		Prefork: true,
		AppName: config.Cfg.App.Name,
	})
	router.Use(middleware.Trace())

	router.Use(cors.New(cors.Config{
		AllowMethods: "GET,POST,PUT,DELETE",
		AllowHeaders: "*",
	}))
	// Logging middleware (optional)
	router.Use(logger.New(logger.Config{
		Format: "[${status}] - ${method} ${path} - ${latency}ms\n",
	}))

	auth.Init(router, db)
	product.Init(router, db)
	transaction.Init(router, db)

	address := fmt.Sprintf(":%v", config.Cfg.App.Port)

	router.Listen(address)
	log.Fatal(router.Listen(address))
}
