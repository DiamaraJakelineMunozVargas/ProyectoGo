package webserver

import (
	//"fmt"
	// "log"

	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
)

func Inicio() {
	app := fiber.New()
	app.Post("/crearcliente", func(c *fiber.Ctx) error {
		type Cliente struct {
			Nombre    string `json:"nombre"`
			MysqlPort int    `json:"mysql_port"`
			MongoPort int    `json:"mongo_port"`
			DicomPort int    `json:"dicom_port"`
			HttpPort  int    `json:"http_port"`
			AppPort   int    `json:"app_port"`
		}
		var cliente Cliente
		err := c.BodyParser(&cliente)
		if err != nil {
			return c.Status(500).SendString("Error al parsear el body")
		}
		crearEntornoCliente(cliente.Nombre)
		plantilla, err := os.ReadFile("compose.template")
		if err != nil {
			return c.Status(500).SendString("No se encontro la plantilla")
		}
		//lo busca y lo reemplaza
		nuevoContent := string(plantilla)
		nuevoContent = strings.ReplaceAll(nuevoContent, "${HOSPITAL_NAME}", cliente.Nombre)
		nuevoContent = strings.ReplaceAll(nuevoContent, "${MYSQL_PORT}", strconv.Itoa(cliente.MysqlPort))
		nuevoContent = strings.ReplaceAll(nuevoContent, "${MONGO_PORT}", strconv.Itoa(cliente.MongoPort))
		nuevoContent = strings.ReplaceAll(nuevoContent, "${DICOM_PORT}", strconv.Itoa(cliente.DicomPort))
		nuevoContent = strings.ReplaceAll(nuevoContent, "${HTTP_PORT}", strconv.Itoa(cliente.HttpPort))
		nuevoContent = strings.ReplaceAll(nuevoContent, "${APP_PORT}", strconv.Itoa(cliente.AppPort))

		// lo guarda
		rutadestino := filepath.Join("/MedicareSoft", cliente.Nombre, "compose.yml")
		err = os.WriteFile(rutadestino, []byte(nuevoContent), 0644)
		if err != nil {
			return c.Status(500).SendString("Error al crear el compose.yml")
		}
		return c.JSON(fiber.Map{"status": "Creado para el cliente " + cliente.Nombre})
	})
	app.Listen(":8080")

}
func crearEntornoCliente(nombre string) {
	carpetas := []string{
		filepath.Join("/Symphony", nombre, "DCM"),
		filepath.Join("/Symphony", nombre, "MYSQL"),
		filepath.Join("/Symphony", nombre, "MONGO"),
		filepath.Join("/Symphony", nombre, "INF"),
		filepath.Join("/Symphony", nombre, "KVSTORE"),
		filepath.Join("/MedicareSoft", nombre, "App"),
	}

	//Crear carpetasddd
	for _, ruta := range carpetas {
		os.MkdirAll(ruta, 0755)

	}
}
