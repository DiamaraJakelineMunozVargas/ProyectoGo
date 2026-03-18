package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func crearEntornoCliente(nombre string) error {
	// 1. Nombrar carpetasss
	carpetas := []string{
		filepath.Join("Symphony", nombre, "DCM"),
		filepath.Join("Symphony", nombre, "MYSQL"),
		filepath.Join("Symphony", nombre, "MONGO"),
		filepath.Join("Symphony", nombre, "INF"),
		filepath.Join("Symphony", nombre, "KVSTORE"),
		filepath.Join("MedicareSoft", nombre, "App"),
	}

	//Crear carpetasddd
	for _, ruta := range carpetas {
		if err := os.MkdirAll(ruta, 0755); err != nil {
			return err
		}
	}

	// Plantilla del Docker Compose
	contenido := fmt.Sprintf(`version: '3.8'
services:
  db-%s:
    image: mysql:8.0
    container_name: mysql-%s
    volumes:
      - ../../Symphony/%s/MYSQL:/var/lib/mysql`, nombre, nombre, nombre)

	// 4. Escribir el archivo
	rutaArchivo := filepath.Join("MedicareSoft", nombre, "compose.yml")
	return os.WriteFile(rutaArchivo, []byte(contenido), 0644)
}

func main() {
	app := fiber.New()

	// Middlewares
	app.Use(logger.New())
	app.Use(recover.New())

	// Ruta 1: Inicio
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Servidor de Gestión MedicareSoft Inicializado 🚀")
	})

	// Ruta 2
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "UP"})
	})

	// NUEVA RUTA: Crear cliente (Ejemplo: localhost:8080/crear/Juan)
	app.Get("/crear/:nombre", func(c *fiber.Ctx) error {
		nombre := c.Params("nombre")

		err := crearEntornoCliente(nombre)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "No se pudo crear el entorno",
			})
		}

		return c.JSON(fiber.Map{
			"mensaje": "¡Entorno de " + nombre + " creado con éxito!",
		})
	})

	// Configuración del puerto000
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Iniciando servidor en el puerto %s...", port)
	log.Fatal(app.Listen(":" + port))
}
