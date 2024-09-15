// Tests for package galaxy functions FindGalaxy and GalaxyChannel
package galaxy

import (
	"testing"

	"star-catalog/database"
)

// TestFindGalaxy calls FindGalaxy with a known ugc_number and checks that
// returns the correct galaxy
func TestFindGalaxy(t *testing.T) {
	db := database.InitDB()

	want := Galaxy{Name: "Milky Way", UgcNumber: "ugc_number1"}
	got, err := FindGalaxy(db, "ugc_number1")

	if err != nil {
		t.Fatalf("FindGalaxy %v\n", err)
	}

	if !ValidateGalaxy(got, want) {
		t.Fatalf(`Should have found galaxy with ugc_number1, is %+v`, got)
	}
}

// TestFindGalaxyMissing calls FindGalaxy with a missing ugc_number and checks that
// returns an error
func TestFindGalaxyMissing(t *testing.T) {
	db := database.InitDB()

	want := "FindGalaxy sql: no rows in result set"

	_, err := FindGalaxy(db, "not_an_ugc_number")

	if err.Error() != want {
		t.Fatalf("FindGalaxy expected error %v, got %v\n", want, err)
	}
}

// TestGalaxyChannel calls GalaxyChannel and checks that it reads the
// contents of the database table
func TestGalaxyChannel(t *testing.T) {
	db := database.InitDB()
	galaxy_channel := make(chan Galaxy)
	error_channel := make(chan error, 1)

	want := []Galaxy{
		{Name: "Milky Way", UgcNumber: "ugc_number1"},
		{Name: "Andromeda", UgcNumber: "ugc_number2"},
		{Name: "", UgcNumber: ""},
	}

	go GalaxyChannel(db, galaxy_channel, error_channel)

	// Check the error channel for errors, otherwise verify output
	select {
	case err, ok := <-error_channel:
		if ok {
			t.Fatalf(`GalaxyChannel: %v`, err)
		}
	default:
		for i := 0; i < 3; i++ {
			got := <-galaxy_channel
			if !ValidateGalaxy(got, want[i]) {
				t.Fatalf(`Item read from channel should be %+v, is %+v`, want[i], got)
			}
		}
	}
}

// TestGalaxyChannelEmpty calls galaxy.GalaxyChannel and checks that it works
// when there are no items in the galaxy table. There is nothing put on either channel.
func TestGalaxyChannelEmpty(t *testing.T) {
	db := database.InitDB()

	database.ClearDB(db)

	galaxy_channel := make(chan Galaxy)
	error_channel := make(chan error, 1)

	go GalaxyChannel(db, galaxy_channel, error_channel)

	var want Galaxy // null Galaxy

	// Check for errors on the error_channel (shouldn't be any)
	for err := range error_channel {
		t.Fatalf(`GalaxyChannel %v\n`, err)
	}

	// There won't be anything on the galaxy_channel
	for got := range galaxy_channel {
		if !ValidateGalaxy(got, want) {
			t.Fatalf(`Any reads should return an empty galaxy, is %v`, got)
		}
	}
}
