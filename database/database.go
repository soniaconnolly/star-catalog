// Package database manages the low level database connection with mysql
// and provides functions InitDB and ClearDB. The connection is automatically closed as needed.
// Database connection details are read from config.yml in the root directory of the project.
// When running tests from a subdirectory, it looks for config.yml in the parent directory.
// The schema is available in schema.sql. Use the mysql utility to read it in. See README.md.
package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
)

// InitDB initializes the database connection and seeds the database with test data.
func InitDB() *sql.DB {
	db := connectDB()
	ClearDB(db)
	seedData(db)

	return db
}

// Connect to a mysql database via the connection information in config.yml in the root directory
func connectDB() *sql.DB {
	findConfigFile()

	cfg := mysql.Config{
		User:      viper.GetString("database.dbuser"),
		Passwd:    viper.GetString("database.dbpassword"),
		Net:       viper.GetString("database.net"),
		Addr:      viper.GetString("database.addr"),
		DBName:    viper.GetString("database.dbname"),
		ParseTime: true,
	}

	// Get a database handle.
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	return db
}

// ClearDB removes all data from the database. Used before seeding the database, and also from tests.
func ClearDB(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM galaxies")
	if err != nil {
		return fmt.Errorf("clearDB: %v", err)
	}

	_, err = db.Exec("DELETE FROM stars")
	if err != nil {
		return fmt.Errorf("clearDB: %v", err)
	}

	return nil
}

// Seed the database with test data
// I looked for exception handling to avoid checking for an error at every step,
// but apparently Go doesn't have that.
func seedData(db *sql.DB) error {
	galaxy_id1, err := addGalaxy(db, "ugc_number1", "Milky Way")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}
	_, err = addStar(db, galaxy_id1, "Sun", "gaia_catalogue_id1")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}
	_, err = addStar(db, galaxy_id1, "Alpha Centauri", "gaia_catalogue_id2")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}

	galaxy_id2, err := addGalaxy(db, "ugc_number2", "Andromeda")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}
	_, err = addStar(db, galaxy_id2, "Star3", "gaia_catalogue_id3")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}
	_, err = addStar(db, galaxy_id2, "Star4", "gaia_catalogue_id4")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}
	_, err = addStar(db, galaxy_id2, "Star5", "gaia_catalogue_id5")
	if err != nil {
		return fmt.Errorf("seedData: %v", err)
	}

	return nil
}

// Add a Galaxy with the given details to the galaxies table
func addGalaxy(db *sql.DB, ugc_number string, name string) (int64, error) {
	result, err := db.Exec("INSERT INTO galaxies (ugc_number, name) VALUES (?, ?)", ugc_number, name)
	if err != nil {
		return 0, fmt.Errorf("addGalaxy: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addGalaxy: %v", err)
	}
	return id, nil
}

// Add a Star with the given details to the stars table
func addStar(db *sql.DB, galaxy_id int64, name string, gaia_catalogue_id string) (int64, error) {
	result, err := db.Exec("INSERT INTO stars (galaxy_id, name, gaia_catalogue_id) VALUES (?, ?, ?)", galaxy_id, name, gaia_catalogue_id)
	if err != nil {
		return 0, fmt.Errorf("addStar: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addStar: %v", err)
	}
	return id, nil
}

// Find config file from tests in a subdirectory, but don't look in the parent directory from main.
// That is: Look in ./ first, and only look in ../ if it's not found. More loop iterations can be
// added for deeper directory structures.
// Taken from https://stackoverflow.com/questions/66683505/handling-viper-config-file-path-during-go-tests
func findConfigFile() {
	var err error
	path := "./"
	for i := 0; i < 2; i++ { // Increase iterations for deeper directory structure
		viper.AddConfigPath(path)
		viper.SetConfigName("config")
		viper.SetConfigType("yml")
		err = viper.ReadInConfig()
		if err != nil {
			if strings.Contains(err.Error(), "Not Found") {
				path = path + "../"
				continue
			}
			panic("panic in config parser : " + err.Error())
		} else {
			break
		}
	}
}
