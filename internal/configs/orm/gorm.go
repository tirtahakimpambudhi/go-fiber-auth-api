package orm

import (
	"fmt"
	sqlconfig "github.com/tirtahakimpambudhi/restful_api/internal/configs/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	// once ensures that the database connection is initialized only once
	once sync.Once
	// db holds the database connection instance
	db *gorm.DB
	// err holds any error encountered during database operations
	err error
)

// NewGorm initializes and returns a singleton instance of a Gorm database connection
func NewGorm() (*gorm.DB, error) {
	// use once.Do to execute the initialization code only once
	once.Do(func() {
		// create a new configuration for SQL
		config, err := sqlconfig.NewConfig()
		if err != nil {
			// if there's an error creating the config, return immediately
			return
		}

		// format the data source name with the configuration settings
		dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta", config.Host, config.User, config.Password, config.Name, config.Port)

		// open a new connection to the database using Gorm
		db, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{
			Logger:                 nil,  // enable SQL logging
			SkipDefaultTransaction: true, // skip default transaction handling
			PrepareStmt:            true, // prepare statements for better performance
		})

		if err != nil {
			// if there's an error opening the database, return immediately
			return
		}

		// retrieve the underlying SQL database instance from Gorm
		sqlDB, err := db.DB()
		if err != nil {
			// if there's an error retrieving the SQL database, return immediately
			return
		}

		// configure database connection pool settings
		sqlDB.SetConnMaxLifetime(30 * time.Minute) // set maximum connection lifetime
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)  // set maximum idle connection time
		sqlDB.SetMaxIdleConns(10)                  // set maximum number of idle connections
		sqlDB.SetMaxOpenConns(100)                 // set maximum number of open connections
	})

	// return the database instance and any error encountered
	return db, err
}
