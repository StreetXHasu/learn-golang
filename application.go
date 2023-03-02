package main

import (
	"app-learn-golang/controllers"
	controllerAdmin "app-learn-golang/controllers/admin"
	"app-learn-golang/initializers"
	"app-learn-golang/middleware"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/expvar"
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

	initializers.ConnectDB(&config)
}

func main() {
	config, _ := initializers.LoadConfig(".") // load environment variables
	port := config.AppPort

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
	log.Printf("Listening on port %s\n\n", port)
	app.Use(cors.New())
	app.Use(expvar.New())
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World!")
	})
	app.Mount("/api", micro)
	micro.Route("/auth", func(router fiber.Router) {
		router.Post("/login", controllers.SignInUser)
		router.Post("/register", controllers.SignUpUser)
		router.Get("/logout", controllers.LogoutUser)
	})

	users := micro.Group("/users", middleware.DeserializeUser)
	users.Get("/me", controllers.GetMe)

	admin := micro.Group("/admin", middleware.DeserializeUser, middleware.Permissions)
	admin.Get("/users", controllerAdmin.UserIndex)

	micro.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "This is API!",
		})
	})
	micro.Get("/version", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"version": "1.1.3",
			"status":  "success",
		})
	})

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
