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
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
)

// Convierte el JSON que llega en un struct de Go.
type Cliente struct {
	gorm.Model        // incerta 4 campos automaticamente :)
	Nombre     string `gorm:"unique;not null" json:"nombre"` //unique para que sea unico
	MysqlPort  int    `json:"mysql_port"`
	MongoPort  int    `json:"mongo_port"`
	DicomPort  int    `json:"dicom_port"`
	HttpPort   int    `json:"http_port"`
	AppPort    int    `json:"app_port"`
}

var db *gorm.DB

func Inicio() {
	var err error
	db, err = gorm.Open(sqlite.Open("Medicaresoft.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database") // se conecta a sqlite
	}

	db.AutoMigrate(&Cliente{}) // se crea la tabla
	app := fiber.New()
	app.Post("/crearcliente", func(c *fiber.Ctx) error {

		type EstructuraNombre struct { // solo pedimos el nombre
			Nombre string `json:"nombre"`
		}

		var ren EstructuraNombre
		if err := c.BodyParser(&ren); err != nil {
			return c.Status(400).SendString("Error al parsear el body")
		}
		var existe Cliente // si ya existe el cliente
		result := db.Where("nombre = ?", ren.Nombre).First(&existe)
		if result.Error == nil {
			return c.Status(400).SendString("Error: Este nombre ya existe")
		}
		// calcular puertos
		nuevoCliente := CalcularPuertos(ren.Nombre)
		// guardar en la BD
		if err := db.Create(&nuevoCliente).Error; err != nil {
			return c.Status(500).SendString("Error al guardar cliente")
		}
		crearEntornoCliente(nuevoCliente.Nombre)
		GenerarArchivoCompose(nuevoCliente)
		return c.JSON(fiber.Map{
			"status":  "OK, ya está en línea, ya está andando",
			"cliente": nuevoCliente,
		})

	})
	// metodo get
	app.Get("/clientes", func(c *fiber.Ctx) error {
		var todos []Cliente
		// busca todos los registros en la tabla :v
		db.Find(&todos)
		return c.JSON(todos)
	})

	app.Listen(":8080")

}
func CalcularPuertos(nombre string) Cliente {
	var ultimo Cliente
	// Buscamos el último cliente para ver sus puertos
	res := db.Order("id desc").First(&ultimo)

	// Puertos base si es el primer cliente
	if res.Error != nil {
		return Cliente{
			Nombre:    nombre,
			MysqlPort: 3306,
			MongoPort: 27017,
			DicomPort: 11112,
			HttpPort:  80,
			AppPort:   3000,
		}
	}

	// Si ya existen clientes, sumamos +1 al último puerto
	return Cliente{
		Nombre:    nombre,
		MysqlPort: ultimo.MysqlPort + 1,
		MongoPort: ultimo.MongoPort + 1,
		DicomPort: ultimo.DicomPort + 1,
		HttpPort:  ultimo.HttpPort + 1,
		AppPort:   ultimo.AppPort + 1,
	}
}
func GenerarArchivoCompose(c Cliente) {
	plantilla, _ := os.ReadFile("compose.template")
	// if err != nil {
	// 	return c.Status(500).SendString("No se encontro la plantilla")
	// }
	//lo busca y lo reemplaza
	nuevoContent := string(plantilla)
	nuevoContent = strings.ReplaceAll(nuevoContent, "${HOSPITAL_NAME}", c.Nombre)
	nuevoContent = strings.ReplaceAll(nuevoContent, "${MYSQL_PORT}", strconv.Itoa(c.MysqlPort))

	nuevoContent = strings.ReplaceAll(nuevoContent, "${MONGO_PORT}", strconv.Itoa(c.MongoPort))
	nuevoContent = strings.ReplaceAll(nuevoContent, "${DICOM_PORT}", strconv.Itoa(c.DicomPort))
	nuevoContent = strings.ReplaceAll(nuevoContent, "${HTTP_PORT}", strconv.Itoa(c.HttpPort))
	nuevoContent = strings.ReplaceAll(nuevoContent, "${APP_PORT}", strconv.Itoa(c.AppPort))

	// lo guarda
	// obtener prefijo
	// prefijo := "/"
	// if os.Getenv("APP_ENV") == "desarrollo" {
	// 	prefijo = ""
	// }
	rutadestino := filepath.Join(obtenerPrefijo(), "MedicareSoft", c.Nombre, "compose.yml")

	//rutadestino := filepath.Join("/MedicareSoft", cliente.Nombre, "compose.yml")
	os.WriteFile(rutadestino, []byte(nuevoContent), 0644)
	// if err != nil {
	// 	return c.Status(500).SendString("Error al crear el compose.yml")
	// }
	// return c.JSON(fiber.Map{"status": "Creado para el cliente " + cliente.Nombre})
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
	// Construimos la ruta al compose.yml del cliente usando tu lógica de prefijo
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

	// se Ejecuta  y capturamos la respuesta de la terminal
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
