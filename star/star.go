// Package star implements the Star struct and function GalaxyStarChannel.
// ValidateStar is availble to use in tests.
// Stars are saved to the stars table.
package star

import (
	"database/sql"
	galaxypkg "star-catalog/galaxy"
	"time"
)

// A Star has an GalaxyId (foreign key to galaxies table), Name, and GaiaCatalogueId
type Star struct {
	Id              int64
	GalaxyId        int64
	Name            string
	GaiaCatalogueId string
	CreatedAt       time.Time
}

// GalaxyStarChannel takes a db connection and a Galaxy, and fills a channel of Star structs
// for the given Galaxy.Id from database table stars
func GalaxyStarChannel(db *sql.DB, galaxy galaxypkg.Galaxy, star_channel chan Star) error {
	// Make sure the channels are closed when the method returns
	defer close(star_channel)

	rows, err := db.Query("SELECT * FROM stars WHERE galaxy_id = ?", galaxy.Id)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var star Star
		if err := rows.Scan(&star.Id, &star.GalaxyId, &star.Name, &star.GaiaCatalogueId, &star.CreatedAt); err != nil {
			return err
		}
		star_channel <- star
	}

	if err := rows.Err(); err != nil {
		return err
	}

	return nil
}

// ValidateStar reports whether the two Stars have the same galaxy_id, ame, and email_address.
// Useful for tests.
func ValidateStar(got Star, want Star) bool {
	return got.GalaxyId == want.GalaxyId && got.Name == want.Name && got.GaiaCatalogueId == want.GaiaCatalogueId
}
