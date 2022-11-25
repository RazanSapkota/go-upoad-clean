package infrastructure

import (
	"example/go-api/lib"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Database modal
type Database struct {
	*gorm.DB
}

// NewDatabase creates a new database instance
func NewDatabase(env *lib.Env) Database {

	username := env.DBUsername
	password := env.DBPassword
	host := env.DBHost
	port := env.DBPort
	dbname := env.DBName

	url := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", username, password, host, port, dbname)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		
	})

	if err != nil {
		log.Println("Url: ", url)
		log.Panic(err)
	}

	log.Println("Database connection established")

	return Database{
		DB: db,
	}
}
