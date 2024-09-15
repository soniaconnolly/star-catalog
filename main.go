// Star-catalog processes all the stars associated with existing galaxies.
// Each galaxy is processed in a separate goroutine.
// Output goes to star-catalog.log
// See README.md for more details.
package main

import (
	"database/sql"
	"log"
	"os"

	"star-catalog/database"
	galaxypkg "star-catalog/galaxy"
	starpkg "star-catalog/star"

	"golang.org/x/sync/errgroup"
)

func main() {
	initLogger()

	// Database handle is passed to methods rather than making it global
	db := database.InitDB()
	Pipeline(db)
}

func Pipeline(db *sql.DB) {
	galaxy_channel := make(chan galaxypkg.Galaxy)
	error_channel := make(chan error, 1)
	var err error

	// Fill the sending channel in a separate goroutine to avoid blocking
	go galaxypkg.GalaxyChannel(db, galaxy_channel, error_channel)

	err = processAllGalaxies(db, galaxy_channel)

	if err != nil {
		log.Printf(`Pipeline %v\n`, err)
	}

	// Check for errors on the error_channel
	for err = range error_channel {
		log.Printf(`Pipeline %v\n`, err)
	}
}

func processAllGalaxies(db *sql.DB, galaxy_channel chan galaxypkg.Galaxy) error {
	// Use an errgroup to wait for all the goroutines to be done and collect any errors
	// https://bostonc.dev/blog/go-errgroup
	g := new(errgroup.Group)

	for galaxy := range galaxy_channel {
		g.Go(func() error {
			return ProcessGalaxy(db, galaxy)
		})
	}

	err := g.Wait()
	return err
}

// ProcessGalaxy takes a database connection and Galaxy, finds the associated stars,
// and calls ProcessStar on each one.
// GalaxyStarChannel is called as a goroutine so that the channel can be processed as
// items are added to it.
func ProcessGalaxy(db *sql.DB, galaxy galaxypkg.Galaxy) error {
	log.Printf("Processing %s galaxy\n", galaxy.UgcNumber)
	var num_stars int
	star_channel := make(chan starpkg.Star)
	error_channel := make(chan error)

	go func() {
		if err := starpkg.GalaxyStarChannel(db, galaxy, star_channel); err != nil {
			error_channel <- err
		}
	}()

	var err error
	var ok bool
	// Check error_channel for errors
	select {
	case err, ok = <-error_channel:
		if ok {
			log.Printf(`ProcessGalaxy %v\n`, err)
			// Identify the galaxy that has been processed, since log output can be interleaved.
			// Not sure we're supposed to log 0 stars processed in case of error?
			log.Printf("%d stars processed for galaxy %s\n", num_stars, galaxy.Name)
			return err
		}
	default:
		// process stars
		for star := range star_channel {
			num_stars++
			// When ProcessStar does real work, would need to add more error handling here
			ProcessStar(galaxy, star)
		}
	}

	// Identify the galaxy that has been processed, since log output can be interleaved.
	log.Printf("%d stars processed for galaxy %s\n", num_stars, galaxy.Name)

	return nil
}

// Process a star, given a Galaxy and a Star.
func ProcessStar(galaxy galaxypkg.Galaxy, star starpkg.Star) {
	log.Printf("Star: %s, Galaxy: %s\n", star.Name, galaxy.Name)
}

// Initialize the logger to go to star-catalog.log and log timestamp
func initLogger() {
	log.SetPrefix("")
	log.SetFlags(log.Ldate | log.Ltime) // Log date and time for ProcessStar

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile("star-catalog.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}
