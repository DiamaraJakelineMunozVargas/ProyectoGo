package webserver

import (
	//"fmt"
	// "log"

	"github.com/gofiber/fiber/v2"
	// "github.com/gofiber/fiber/v2/middleware/logger"
	// "github.com/gofiber/fiber/v2/middleware/recover"
)

func Inicio() {
	ConexionDB() //se llama a la funcion donde esta la base de datos
	app := fiber.New()
	app.Post("/crearcliente", CrearCliente)
	// metodo get
	app.Get("/clientes", Clientes)

	app.Listen(":8080")

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
