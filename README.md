# Star Catalog in Go
## Assignment
Process stars within galaxies stored in a relational database system (RDS) using Go channels and goroutines for parallel processing.

### Schema
```
CREATE TABLE galaxies(
    id INT,
    ugc_number VARCHAR(200),
    name VARCHAR(200),
    created_at TIMESTAMP
);

CREATE TABLE stars(
    id INT,
    galaxy_id INT,
    name VARCHAR(200),
    gaia_catalogue_id VARCHAR(200), 
    created_at TIMESTAMP
);
```
 `ugc_number` is from [Uppsala General Catalog of Galaxies](https://heasarc.gsfc.nasa.gov/W3Browse/galaxy-catalog/ugc.html)
 
`gaia_catalogue_id` is from [Gaia Catalogue of Nearby Stars](https://www.cosmos.esa.int/web/gaia/edr3-gcns)

### Tasks
 1. Function to Retrieve Galaxies: Write a function using the Go SQL package of your choice that returns a channel with the data of the galaxies.
 2. Function to Retrieve Stars: Write a function using the Go SQL package of your choice that receives the data of a galaxy and returns a channel with the stars that belong to that galaxy.
 3. Tests for Functions: Write tests for each of the above functions.
 4. Pipeline Function: Write a function that builds a "pipeline" using the previous two functions to call a fictional processStar() function for each star.
   - processStar() should log the time, star name, and galaxy name using the Log package of your choice.
   - The "pipeline" should process each galaxy in a separate goroutine.
   - The pipeline should log the progress: "Processing <ugc_number> galaxy" when it starts processing a galaxy and "<number_of_stars> stars processed" when it finishes processing them.
 5. Test for Pipeline: Write a test for the pipeline function.
 6. Configuration Handling: The application should read the database connection from a file. Use the configuration package/method of your choice.
 7. Main Function: Write a main() function that:
   • Initializes the database connection.
   • Call the pipeline.

### Requirements
 - Language: Go (Golang)
 - Database Access: Use any Go SQL package you prefer (e.g., database/sql, gorm).
 - Logging: Use any logging package you prefer.
 - Error Handling: Handle errors as you would in a production application.
 - Testing: Ensure that the tests pass and have good coverage. The code itself does not need to be executable, but the tests should demonstrate the functionality.

## Setup

### Initialize database
All these commands assume the current working directory is the top level directory of the project.

Run mysql:
```
$ mysql -u root
```
Inside mysql:
```
create database star_catalog;
use star_catalog;
source ./database/schema.sql
```

Put your local mysql settings into `./config.yml`

### Run tests
```
make test
make test/cover
```
The code segments that aren't covered are checking for database errors. This could be tested by using dependency injection with the DB and using a mock to return the desired errors.

### Run the app
```
make run
```
In a separate window
```
tail -f star-catalog.log
```

## Directories and files
I didn't find a unified best practice for structuring the files of a Go app. Based on this article, I chose a simple package structure separating low level database code, galaxy code, and star code.
https://www.calhoun.io/using-mvc-to-structure-go-web-applications/ 

## Documentation and Tutorials
This is my first Go program, based on reading the following tutorials. It was done as part of an interview process: an evening to read tutorials, two days of coding, and another evening to clean up the code and add documentation. 

- https://go.dev/doc/tutorial/getting-started
- https://go.dev/doc/tutorial/create-module
- https://go.dev/tour/
- https://go.dev/doc/effective_go
- https://go.dev/doc/tutorial/database-access
- https://go.dev/blog/pipelines
- https://go.dev/doc/comment
- https://blog.logrocket.com/handling-go-configuration-viper/
- https://www.sohamkamani.com/golang/time/
- https://go.dev/doc/code
- https://go.dev/talks/2012/concurrency.slide#1
- https://medium.com/@jomzsg/the-easy-way-to-handle-configuration-file-in-golang-using-viper-6b3c88d2ee79

## Notes
### Closing DB connections
Looks like the DB connection is automatically closed in Go, based on example code and this stack overflow response.
https://stackoverflow.com/questions/40587008/how-do-i-handle-opening-closing-db-connection-in-a-go-app

### Error handling
Looks like Go doesn't have exception handling. I looked up error handling in an effort to avoid checking for errors after every line of database.seedData(), but left it as is for now.
https://blog.logrocket.com/error-handling-golang-best-practices/ 

Channel error handling taken from here. 
https://blog.poespas.me/posts/2024/04/29/go-errgroup-and-channel-best-practices/

### Internationalization
All messages shown to the user, including error messages, should be internationalizable. I looked it up, but didn't add internationalized strings since this seems to be a backend app rather than a user-facing one.
https://phrase.com/blog/posts/internationalisation-in-go-with-go-i18n/

### Mocks in tests
In a larger project, using the database in tests might be unacceptably slow, and mocks would be used instead. I looked up mocks in tests, but chose not to use them for this small project.
https://blog.logrocket.com/exploring-go-mocking-methods-gomock-framework/

More generally, I based my test structure on
https://www.digitalocean.com/community/tutorials/how-to-write-unit-tests-in-go-using-go-test-and-the-testing-package
