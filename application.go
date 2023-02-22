package main

import (
	"app-learn-golang/controllers"
	"app-learn-golang/initializers"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"log"
	"os"
)

func init() {
	log.Println("Init application...")
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatalln("Failed to load environment variables! \n", err.Error())
	}
	log.Println("port ", config.DBPort)

	//initializers.ConnectDB(&config)
}

func main() {
	config, _ := initializers.LoadConfig(".") // load environment variables
	port := os.Getenv("PORT")
	if port == "" {
		port = config.AppPort
	}

	f, _ := os.Create("/var/log/golang/golang-server.log")
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println(err)
		}
	}(f)
	log.SetOutput(f)

	app := fiber.New()
	micro := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})
	app.Mount("/api", micro)
	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/login", controllers.SignInUser)
		router.Get("/logout", controllers.LogoutUser)
	})
	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("This is API!")
	})

	app.Get("/api/version", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"version": "1.1.0",
			"status":  "Ok",
		})
	})

	log.Printf("Listening on port %s\n\n", port)
	micro.Get("/healthchecker", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "JWT Authentication with Golang, Fiber",
		})
	})

	micro.All("*", func(c *fiber.Ctx) error {
		path := c.Path()
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "fail",
			"message": fmt.Sprintf("Path: %v does not exists on this server", path),
		})
	})

	log.Fatal(app.Listen(":" + port))
}
