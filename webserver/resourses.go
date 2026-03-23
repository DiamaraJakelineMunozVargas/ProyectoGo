package webserver

import "github.com/gofiber/fiber/v2"

func CrearCliente(c *fiber.Ctx) error {
	type EstructuraNombre struct { // solo pedimos el nombre
		Nombre string `json:"nombre"`
	}

	var ren EstructuraNombre
	if err := c.BodyParser(&ren); err != nil {
		return c.Status(400).SendString("Error al parsear el body")
	}
	var existe Cliente // si ya existe el cliente
	result := DB.Where("nombre = ?", ren.Nombre).First(&existe)
	if result.Error == nil {
		return c.Status(400).SendString("Error: Este nombre ya existe")
	}
	// calcular puertos
	nuevoCliente := CalcularPuertos(ren.Nombre)
	// guardar en la BD
	if err := DB.Create(&nuevoCliente).Error; err != nil {
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
	var todos []Cliente
	// busca todos los registros en la tabla :v
	DB.Find(&todos)
	return c.JSON(todos)
}
