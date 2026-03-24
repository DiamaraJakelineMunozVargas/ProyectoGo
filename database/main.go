package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

var DB *gorm.DB

func ConexionDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("Medicaresoft.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database") // se conecta a sqlite
	}

	DB.AutoMigrate(&Cliente{}) // se crea la tabla
}
