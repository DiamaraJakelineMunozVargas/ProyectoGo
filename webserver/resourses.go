package webserver

import (
	"joseluis244/ProyectoGo/database"

	"github.com/gofiber/fiber/v2"
)

func CrearCliente(c *fiber.Ctx) error {
	type EstructuraNombre struct { // solo pedimos el nombre
		Nombre string `json:"nombre"`
	}

	var ren EstructuraNombre
	if err := c.BodyParser(&ren); err != nil {
		return c.Status(400).SendString("Error al parsear el body")
	}
	var existe database.Cliente // si ya existe el cliente
	result := database.DB.Where("nombre = ?", ren.Nombre).First(&existe)
	if result.Error == nil {
		return c.Status(400).SendString("Error: Este nombre ya existe")
	}
	// calcular puertos
	nuevoCliente := CalcularPuertos(ren.Nombre)
	// guardar en la BD
	if err := database.DB.Create(&nuevoCliente).Error; err != nil {
		return c.Status(500).SendString("Error al guardar cliente")
	}
	CrearEntornoCliente(nuevoCliente.Nombre)
	GenerarArchivoCompose(nuevoCliente)
	return c.JSON(fiber.Map{
		"status":  "OK, ya está en línea, ya está andando",
		"cliente": nuevoCliente,
	})
}

func Clientes(c *fiber.Ctx) error {
	var todos []database.Cliente
	// busca todos los registros en la tabla :v
	database.DB.Find(&todos)
	return c.JSON(todos)
}
func ControlarDocker(c *fiber.Ctx) error {
	nombre := c.Params("nombre")
	accion := c.Params("accion")
	output, err := GestionarDocker(accion, nombre)
	if err != nil {
		return c.Status(500).SendString("Error Docker: " + err.Error() + output)
	}
	return c.SendString("Operacion " + accion + "exitosa para " + nombre)
}
func VerLog(c *fiber.Ctx) error {
	nombre := c.Params("nombre")
	servicio := c.Params("servicio")
	resultado, err := GestionarDocker(nombre, servicio)
	if err != nil {
		return c.Status(500).SendString("Error al ver el log" + resultado)
	}
	return c.SendString(resultado)
}
