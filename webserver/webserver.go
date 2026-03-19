package webserver

import (
	//"fmt"
	// "log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
)

func Inicio() {
	app := fiber.New()
	app.Post("/crearcliente/:nombre", func(c *fiber.Ctx) error {
		nombre := c.Params("nombre")
		crearEntornoCliente(nombre)
		plantilla, err := os.ReadFile("compose.template")
		if err != nil {
			return c.Status(500).SendString("No se encontro la plantilla")
		}
		//lo busca y lo reemplaza
		nuevoContent := string(plantilla)
		nuevoContent = strings.ReplaceAll(nuevoContent, "${HOSPITAL_NAME}", nombre)

		// nuevoContent = strings.ReplaceAll(nuevoContent, "${MYSQL_PORT}", "3306")
		// nuevoContent = strings.ReplaceAll(nuevoContent, "${MONGO_PORT}", "27017")
		// nuevoContent = strings.ReplaceAll(nuevoContent, "${DICOM_PORT}", "3101")
		// nuevoContent = strings.ReplaceAll(nuevoContent, "${HTTP_PORT}", "8042")
		// nuevoContent = strings.ReplaceAll(nuevoContent, "${APP_PORT}", "4000")

		// lo guarda
		rutadestino := filepath.Join("MedicareSoft", nombre, "compose.yml")
		err = os.WriteFile(rutadestino, []byte(nuevoContent), 0644)
		if err != nil {
			return c.Status(500).SendString("Error al crear el compose.yml")
		}
		return c.JSON(fiber.Map{"status": "Creado para el cliente " + nombre})
	})
	app.Listen(":8080")

}
func crearEntornoCliente(nombre string) {
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
		os.MkdirAll(ruta, 0755)

	}
}
