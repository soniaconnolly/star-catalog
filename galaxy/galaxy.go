// Package galaxy implements the Galaxy struct, and GalaxyChannel and FindGalaxy functions.
// ValidateGalaxy is available to use in tests.
// Galaxies are saved to the galaxies table.
package galaxy

import (
	"database/sql"
	"fmt"
	"time"
)

// A Galaxy has a UGC number, and has Stars associated with it
type Galaxy struct {
	Id        int64
	Name      string
	UgcNumber string
	CreatedAt time.Time
}

// GalaxyChannel takes a database connection, and fills galaxy_channel with Galaxy structs
// from database table galaxies.
// If there is an error it is put on error_channel.
func GalaxyChannel(db *sql.DB, galaxy_channel chan Galaxy, error_channel chan error) {
	// Make sure the channels are closed when the method returns
	defer close(galaxy_channel)
	defer close(error_channel)
	rows, err := db.Query("SELECT * FROM galaxies")
	if err != nil {
		error_channel <- err
		return
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var galaxy Galaxy
		if err := rows.Scan(&galaxy.Id, &galaxy.Name, &galaxy.UgcNumber, &galaxy.CreatedAt); err != nil {
			error_channel <- err
			return
		}
		galaxy_channel <- galaxy
	}

	if err := rows.Err(); err != nil {
		error_channel <- err
	}
}

// FindGalaxy takes a database connection and an UgcNumber, and returns the Galaxy struct
// found, or an error if not.
func FindGalaxy(db *sql.DB, ugc_number string) (Galaxy, error) {
	var galaxy Galaxy

	err := db.QueryRow("SELECT * FROM galaxies WHERE ugc_number = ?", ugc_number).
		Scan(&galaxy.Id, &galaxy.Name, &galaxy.UgcNumber, &galaxy.CreatedAt)

	if err != nil {
		return galaxy, fmt.Errorf("FindGalaxy %v", err)
	}

	return galaxy, nil
}

// ValidateGalaxy reports whether the two Galaxies have the same ugc_number.
// Useful for tests.
func ValidateGalaxy(got Galaxy, want Galaxy) bool {
	return got.UgcNumber == want.UgcNumber
}
