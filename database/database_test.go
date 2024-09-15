// Tests for InitDB and ClearDB
package database

import (
	"testing"
)

// TestInitDB calls database.InitDB and checks that the database
// is initialized
func TestInitDB(t *testing.T) {
	db := InitDB()

	var got int
	want := 2

	err := db.QueryRow("SELECT COUNT(*) FROM galaxies").Scan(&got)
	if got != want || err != nil {
		t.Fatalf(`Galaxy rows should be %d, is %d, %v`, want, got, err)
	}

	want = 5

	err = db.QueryRow("SELECT COUNT(*) FROM stars").Scan(&got)
	if got != want || err != nil {
		t.Fatalf(`Star rows should be %d, is %d, %v`, want, got, err)
	}
}

// TestClearDB calls database.ClearDB and confirms that there are no
// items in galaxies and stars
func TestClearDB(t *testing.T) {
	db := InitDB()
	ClearDB(db)

	var got int
	want := 0

	err := db.QueryRow("SELECT COUNT(*) FROM galaxies").Scan(&got)
	if got != want || err != nil {
		t.Fatalf(`Galaxy rows should be %d, is %d, %v`, want, got, err)
	}

	err = db.QueryRow("SELECT COUNT(*) FROM stars").Scan(&got)
	if got != want || err != nil {
		t.Fatalf(`Star rows should be %d, is %d, %v`, want, got, err)
	}
}
