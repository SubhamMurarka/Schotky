package main

import (
	"fmt"
	"log"
	"os"

	"github.com/SubhamMurarka/Schotky/Config"
	"github.com/SubhamMurarka/Schotky/Dynamo"
	handlers "github.com/SubhamMurarka/Schotky/Handlers"
	services "github.com/SubhamMurarka/Schotky/Services"
	zookeepercounter "github.com/SubhamMurarka/Schotky/ZookeeperCounter"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func setupRoutes(app *fiber.App, h *handlers.Handler) {
	app.Get("/:url", h.ResolveUrl)
	app.Post("/api/v1", h.ShortenUrl)
}

func main() {
	//initialising zookeeper
	dir, _ := os.Getwd()
	fmt.Println(dir)

	// initialise Dynamo and Dax
	DDclient := Dynamo.NewDynamoDaxClient()
	DDclient.CreateTable()
	DDclient.EnableTTL()

	// var client zookeepercounter.ZooKeeperClient
	client := zookeepercounter.NewZooKeeperClient()
	err := client.Connect()
	if err != nil {
		log.Fatal("error connecting to zookeeper, ", err)
	}
	client.CreatePersistentNodes()
	defer client.Close()

	//initialising services
	ss := services.NewShortenServiceObj(DDclient, client)
	rs := services.NewResolveServiceObj(DDclient)

	//initialising handler
	h := handlers.NewHandler(ss, rs)

	app := fiber.New()

	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // for development, specify the exact origins in production
		AllowMethods: "GET, POST, HEAD, PUT, DELETE, PATCH",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	setupRoutes(app, h)

	APP_PORT := Config.Cfg.APP_PORT

	log.Fatal(app.Listen(APP_PORT))
}
