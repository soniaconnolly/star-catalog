// Tests for Pipeline, ProcessGalaxy, and ProcessStar
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"star-catalog/database"
	"star-catalog/galaxy"
	galaxypkg "star-catalog/galaxy"
	starpkg "star-catalog/star"
	"strings"
	"testing"

	// "database"
	// galaxypkg "galaxy"
	// starpkg "star"

	"github.com/stretchr/testify/assert"
)

// TestPipelineOutput calls Pipeline and checks the log output with logging redirected to a pipe
// This confirms that the correct lines were received, possibly interleaved.
func TestPipelineOutput(t *testing.T) {
	scanner, reader, writer := mockLogger(t) // turn this off when debugging or developing as you will miss output!
	defer resetLogger(reader, writer)

	// Log output is listed in order and then sorted because messages can arrive interleaved
	want := []string{
		"Processing ugc_number1 galaxy",
		"Star: Sun, Galaxy: Milky Way",
		"Star: Alpha Centauri, Galaxy: Milky Way",
		"2 stars processed for galaxy Milky Way",
		"Processing ugc_number2 galaxy",
		"Star: Star3, Galaxy: Andromeda",
		"Star: Star4, Galaxy: Andromeda",
		"Star: Star5, Galaxy: Andromeda",
		"3 stars processed for galaxy Andromeda",
	}
	sort.Strings(want)

	var got []string

	db := database.InitDB()
	go Pipeline(db)

	// Log lines can be interleaved. Collect the lines, and sort to compare with expected array.
	for i := 0; i < 9; i++ {
		scanner.Scan()         // blocks until a new line is written to the pipe
		line := scanner.Text() // the last line written to the scanner
		got = append(got, line)
	}

	sort.Strings(got)

	for i := 0; i < 9; i++ {
		if !strings.Contains(got[i], want[i]) {
			t.Fatalf(`Pipeline log should match %s, is %s`, want[i], got)
		}
	}
}

func TestProcessGalaxyOutput(t *testing.T) {
	scanner, reader, writer := mockLogger(t) // turn this off when debugging or developing as you will miss output!
	defer resetLogger(reader, writer)

	want := []string{
		"Processing ugc_number1 galaxy",
		"Star: Sun, Galaxy: Milky Way",
		"Star: Alpha Centauri, Galaxy: Milky Way",
		"2 stars processed for galaxy Milky Way",
	}

	db := database.InitDB()

	mbox, err := galaxy.FindGalaxy(db, "ugc_number1")
	if err != nil {
		t.Fatalf(`ProcessGalaxy %v\n`, err)
	}

	err = ProcessGalaxy(db, mbox)

	if err != nil {
		t.Fatalf(`ProcessGalaxy %v\n`, err)
	}

	// Verify the log output
	for i := 0; i < 4; i++ {
		scanner.Scan()        // blocks until a new line is written to the pipe
		got := scanner.Text() // the last line written to the scanner

		if !strings.Contains(got, want[i]) {
			t.Fatalf(`ProcessGalaxy log should match %s, is %s`, want[i], got)
		}
	}
}

// TestProcessStarOutput calls ProcessStar and checks the log output with logging redirected to a pipe
// This confirms that the logs have the required format including time, name, and galaxy token
func TestProcessStarOutput(t *testing.T) {
	scanner, reader, writer := mockLogger(t) // turn this off when debugging or developing as you will miss output!
	defer resetLogger(reader, writer)

	// Go doesn't seem to have a good way to mock time.now(), so just check that a date and time
	// is logged, since the functionality is built into the log package.
	// Output format: 2024/08/11 15:01:47 Star: star_name, Galaxy: galaxy_name
	want := `\d\d\d\d\/\d\d\/\d\d \d\d\:\d\d\:\d\d Star: Sun, Galaxy: Milky Way`

	mbox := galaxypkg.Galaxy{UgcNumber: "ugc_number1", Name: "Milky Way"}
	usr := starpkg.Star{Name: "Sun", GaiaCatalogueId: "gaia_catalogue_id1"}

	ProcessStar(mbox, usr)

	scanner.Scan()        // blocks until a new line is written to the pipe
	got := scanner.Text() // the last line written to the scanner

	matched, err := regexp.MatchString(want, got)

	if !matched || err != nil {
		t.Fatalf(`ProcessStar log should match %s, is %s`, want, got)
	}
}

// Taken from https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
func mockLogger(t *testing.T) (*bufio.Scanner, *os.File, *os.File) {
	reader, writer, err := os.Pipe()
	if err != nil {
		assert.Fail(t, "couldn't get os Pipe: %v", err)
	}
	log.SetOutput(writer)

	return bufio.NewScanner(reader), reader, writer
}

func resetLogger(reader *os.File, writer *os.File) {
	err := reader.Close()
	if err != nil {
		fmt.Println("error closing reader was ", err)
	}
	if err = writer.Close(); err != nil {
		fmt.Println("error closing writer was ", err)
	}
	log.SetOutput(os.Stderr)
}
