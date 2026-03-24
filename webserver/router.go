package webserver

import "github.com/gofiber/fiber/v2"

func Router(router fiber.Router) {
	router.Post("/crearcliente", CrearCliente)
	// metodo get
	router.Get("/clientes", Clientes)

	router.Post("/docker/:nombre/:accion", ControlarDocker)
	router.Get("/log/:nombre/:servicio", VerLog)
}
