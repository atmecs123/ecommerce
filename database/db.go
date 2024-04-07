package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const DB_USERNAME = "root"
const DB_PASSWORD = "root"
const DB_NAME = "test"
const DB_HOST = "localhost"
const DB_PORT = "3306"

var DB *gorm.DB

func InitDb() (*gorm.DB, error) {
	var err error
	DB, err = ConnectDB()
	return DB, err
}

func ConnectDB() (*gorm.DB, error) {
	dsn := DB_USERNAME + ":" + DB_PASSWORD + "@tcp" + "(" + DB_HOST + ":" + DB_PORT + ")/" + DB_NAME + "?" + "parseTime=true&loc=Local"
	fmt.Println("Full database connection name", dsn)
	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Error connecting to my sql database", err)
		return gormDB, err
	}

	return gormDB, err

}
