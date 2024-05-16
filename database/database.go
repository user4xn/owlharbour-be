package database

import (
	"fmt"
	"log"
	"owlharbour-api/pkg/util"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbConn *gorm.DB
	once   sync.Once
)

func CreateConnection() {
	// Create database configuration information
	conf := dbConfig{
		User: util.GetEnv("PG_USER", "postgres"),      // Update to your PostgreSQL username
		Pass: util.GetEnv("PG_PASS", ""),              // Update to your PostgreSQL password
		Host: util.GetEnv("PG_HOST", "localhost"),     // Update to your PostgreSQL host
		Port: util.GetEnv("PG_PORT", "5432"),          // Update to your PostgreSQL port
		Name: util.GetEnv("PG_DB_NAME", "simpel_api"), // Update to your PostgreSQL database name
	}

	// PostgreSQL connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", conf.Host, conf.Port, conf.User, conf.Pass, conf.Name)

	var err error
	dbConn, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to PostgreSQL: %v", err)
	}
}

func GetConnection() *gorm.DB {
	// Check db connection, if exist return the memory address of the db connection
	if dbConn == nil {
		CreateConnection()
	}
	return dbConn
}
