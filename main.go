package main

import (
	"joseluis244/ProyectoGo/database"
	"joseluis244/ProyectoGo/webserver"
)

func main() {
	database.ConexionDB()
	webserver.Inicio()
}
