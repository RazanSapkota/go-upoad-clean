package initialize

import (
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

  var DB *gorm.DB
func ConnectDB(){
	var err error
	
		// refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn := os.Getenv("DB_URL")
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			QueryFields: true,
		})
	  
		if err!=nil {
			log.Fatal("Failed to connect to database")
		}
}