package webserver

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"joseluis244/ProyectoGo/database"
)

func CalcularPuertos(nombre string) database.Cliente {
	var ultimo database.Cliente
	// Buscamos el último cliente para ver sus puertos
	res := database.DB.Order("id desc").First(&ultimo)

	// Puertos base si es el primer cliente
	if res.Error != nil {
		return database.Cliente{
			Nombre:    nombre,
			MysqlPort: 3306,
			MongoPort: 27017,
			DicomPort: 11112,
			HttpPort:  80,
			AppPort:   3000,
		}
	}

	// Si ya existen clientes, sumamos +1 al último puerto
	return database.Cliente{
		Nombre:    nombre,
		MysqlPort: ultimo.MysqlPort + 1,
		MongoPort: ultimo.MongoPort + 1,
		DicomPort: ultimo.DicomPort + 1,
		HttpPort:  ultimo.HttpPort + 1,
		AppPort:   ultimo.AppPort + 1,
	}
}
func GenerarArchivoCompose(c database.Cliente) {
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
	rutadestino := filepath.Join(ObtenerPrefijo(), "MedicareSoft", c.Nombre, "compose.yml")

	//rutadestino := filepath.Join("/MedicareSoft", cliente.Nombre, "compose.yml")
	os.WriteFile(rutadestino, []byte(nuevoContent), 0644)
	// if err != nil {
	// 	return c.Status(500).SendString("Error al crear el compose.yml")
	// }
	// return c.JSON(fiber.Map{"status": "Creado para el cliente " + cliente.Nombre})
}
func CrearEntornoCliente(nombre string) {

	carpetas := []string{
		filepath.Join(ObtenerPrefijo(), "Symphony", nombre, "DCM"),
		filepath.Join(ObtenerPrefijo(), "Symphony", nombre, "MYSQL"),
		filepath.Join(ObtenerPrefijo(), "Symphony", nombre, "MONGO"),
		filepath.Join(ObtenerPrefijo(), "Symphony", nombre, "INF"),
		filepath.Join(ObtenerPrefijo(), "Symphony", nombre, "KVSTORE"),
		filepath.Join(ObtenerPrefijo(), "MedicareSoft", nombre, "App"),
	}

	for _, ruta := range carpetas {
		os.MkdirAll(ruta, 0755)
	}
}
func ObtenerPrefijo() string {
	// Si la variable no existe o esta en desarrollo, devolvemos vacío (carpeta local)
	if os.Getenv("APP_ENV") == "produccion" {
		return "/"
	}
	return ""
}
func GestionarDocker(accion string, nombre string, servicio ...string) (string, error) {
	// Construimos la ruta al compose.yml del cliente usando tu lógica de prefijo
	rutaCompose := filepath.Join(ObtenerPrefijo(), "MedicareSoft", nombre, "compose.yml")

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
		ContainerName := nombre + "_" + servicio[0]
		cmd = exec.Command("docker", "logs", "--tail", "100", ContainerName)
	default:
		return "Acción no permitida", nil
	}

	// se Ejecuta  y capturamos la respuesta de la terminal
	out, err := cmd.CombinedOutput()
	return string(out), err
}
