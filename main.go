package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	//Crer la app de Fiber
	/*app := fiber.New()
	app.Use(LogginMiddleware)

	log.Fatal(app.Listen(":8080"))*/
	// Crear la instancia de Fiber v2
	app := fiber.New()

	// Middlewares

	app.Use(logger.New())  // Para ver quién entra en la terminal
	app.Use(recover.New()) // Para que el servidor no se caiga si hay un error

	// Ruta de prueba
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hola mundo")
	})

	// segunda Ruta revision
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP"})
	})

	//Configuración del puerto
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Encender el servidor
	log.Printf("Iniciando servidor en el puerto %s...", port)
	log.Fatal(app.Listen(":" + port))
}

// func LogginMiddleware(c *fiber.Ctx) error {
// 	log.Printf("Solicitud recibida: %s %s ", c.Method(), c.Path())
// 	return c.Next()
// }
