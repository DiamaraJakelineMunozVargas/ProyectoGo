package webserver

import (
	//"fmt"
	// "log"

	"os"
	"os/exec"
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
		//Convierte el JSON que llega en un struct de Go.
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
		// obtener prefijo
		// prefijo := "/"
		// if os.Getenv("APP_ENV") == "desarrollo" {
		// 	prefijo = ""
		// }
		rutadestino := filepath.Join(obtenerPrefijo(), "MedicareSoft", cliente.Nombre, "compose.yml")

		//rutadestino := filepath.Join("/MedicareSoft", cliente.Nombre, "compose.yml")
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
		filepath.Join(obtenerPrefijo(), "Symphony", nombre, "DCM"),
		filepath.Join(obtenerPrefijo(), "Symphony", nombre, "MYSQL"),
		filepath.Join(obtenerPrefijo(), "Symphony", nombre, "MONGO"),
		filepath.Join(obtenerPrefijo(), "Symphony", nombre, "INF"),
		filepath.Join(obtenerPrefijo(), "Symphony", nombre, "KVSTORE"),
		filepath.Join(obtenerPrefijo(), "MedicareSoft", nombre, "App"),
	}

	for _, ruta := range carpetas {
		os.MkdirAll(ruta, 0755)
	}
}
func obtenerPrefijo() string {
	// Si la variable no existe o esta en desarrollo, devolvemos vacío (carpeta local)
	if os.Getenv("APP_ENV") == "produccion" {
		return "/"
	}
	return ""
}
func gestionarDocker(accion string, nombre string) (string, error) {
	// 1. Construimos la ruta al compose.yml del cliente usando tu lógica de prefijo
	rutaCompose := filepath.Join(obtenerPrefijo(), "MedicareSoft", nombre, "compose.yml")

	// Preparamos el comando según la acción
	// Ejemplo: docker compose -f MedicareSoft/Juan/compose.yml up -d
	var cmd *exec.Cmd
	switch accion {
	case "start":
		cmd = exec.Command("docker", "compose", "-f", rutaCompose, "up", "-d")
	case "stop":
		cmd = exec.Command("docker", "compose", "-f", rutaCompose, "stop")
	case "restart":
		cmd = exec.Command("docker", "compose", "-f", rutaCompose, "restart")
	case "logs":
		cmd = exec.Command("docker", "compose", "-f", rutaCompose, "logs", "--tail=20")
	default:
		return "Acción no permitida", nil
	}

	// se Ejecuta  y captura la respuesta de la terminal
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// comandos para la terminal
// APP_ENV=desarrollo go run main.go
// % curl -X POST http://localhost:8080/crearcliente \
// -H "Content-Type: application/json" \
// -d '{
//   "nombre": "Jorge",
//   "mysql_port": 3307,
//   "mongo_port": 27018,
//   "dicom_port": 3102,
//   "http_port": 8043,
//   "app_port": 4001}'

// func crearEntornoCliente(nombre string) {
// 	carpetas := []string{
// 		filepath.Join("/Symphony", nombre, "DCM"),
// 		filepath.Join("/Symphony", nombre, "MYSQL"),
// 		filepath.Join("/Symphony", nombre, "MONGO"),
// 		filepath.Join("/Symphony", nombre, "INF"),
// 		filepath.Join("/Symphony", nombre, "KVSTORE"),
// 		filepath.Join("/MedicareSoft", nombre, "App"),
// 	}

// 	//Crear carpetasddd
// 	for _, ruta := range carpetas {
// 		os.MkdirAll(ruta, 0755)

// 	}
// }
