// Tests for GalaxyStarChannel
package star

import (
	"fmt"
	"testing"

	"star-catalog/database"
	"star-catalog/galaxy"
)

// TestGalaxyStarChannel calls star.GalaxyStarChannel and checks that it finds all the stars
// for the given galaxy id
func TestGalaxyStarChannel(t *testing.T) {
	db := database.InitDB()
	star_channel := make(chan Star)

	galaxy1, err := galaxy.FindGalaxy(db, "ugc_number1")
	if err != nil {
		t.Fatalf("GalaxyStarChannel: %v", err)
	}

	want := []Star{
		{GalaxyId: galaxy1.Id, Name: "Sun", GaiaCatalogueId: "gaia_catalogue_id1"},
		{GalaxyId: galaxy1.Id, Name: "Alpha Centauri", GaiaCatalogueId: "gaia_catalogue_id2"},
		{GalaxyId: 0, Name: "", GaiaCatalogueId: ""},
	}

	error_channel := make(chan error)

	go func() {
		if err := GalaxyStarChannel(db, galaxy1, star_channel); err != nil {
			error_channel <- err
		}
	}()

	var got Star

	for i := 0; i < 3; i++ {
		// Check for errors on the error_channel and otherwise verify output
		select {
		case err = <-error_channel:
			fmt.Printf("GalaxyStarChannel %v\n", err)
		case got = <-star_channel:
			if !ValidateStar(got, want[i]) {
				t.Fatalf(`Item read from channel should be %+v, is %+v`, want[i], got)
			}
		}
	}
}

// TestGalaxyStarChannelEmpty calls star.GalaxyStarChannel and checks that it works
// when there are no items in the stars table for the given id
func TestGalaxyStarChannelEmpty(t *testing.T) {
	db := database.InitDB()

	database.ClearDB(db)

	galaxy := galaxy.Galaxy{Id: 1, UgcNumber: "ugc_number1"}
	star_channel := make(chan Star)
	err := GalaxyStarChannel(db, galaxy, star_channel)

	var want Star // null Star

	if err != nil {
		t.Fatalf(`GalaxyStarChannel %v\n`, err)
	}

	// There won't be anything on the star_channel
	for got := range star_channel {
		if !ValidateStar(got, want) {
			t.Fatalf(`Any reads should return an empty star, is %v`, got)
		}
	}
}
