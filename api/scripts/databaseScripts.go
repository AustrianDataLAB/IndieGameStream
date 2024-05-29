package scripts

import (
	"api/shared"
	"database/sql"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

func ConnectToDatabase() *sql.DB {
	// connect to db using standard Go database/sql API
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		os.Getenv("MYSQL_ROOT_USER"),
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		//println(connectionString)
		log.Fatal(err)
	}

	return db
}

// MigrateDatabase applies the migration scripts from folder migrations,
// if they have not been applied already.
func MigrateDatabase(db *sql.DB) {
	log.Println("Starting migrations...")

	//Get the migrations that has been applied
	migrations := getMigrationIds(db)

	//Load all migration scripts
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Fatal(err)
	}

	//Sort the list of fileNames
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	sort.Strings(fileNames)

	//For each migration script
	for _, fileName := range fileNames {
		//Check if its a valid filename
		match, err := regexp.MatchString(`\d*_.*[.]sql`, fileName)
		if err != nil {
			log.Fatal(err)
		}
		if !match {
			continue
		}

		//If the file is not a .sql file, ignore it
		if fileName[len(fileName)-4:] != ".sql" {
			continue
		}

		//get its id
		migrationId, err := strconv.Atoi(strings.Split(fileName, "_")[0])
		if err != nil {
			log.Fatal(err)
		}
		//If the migration has not been applied yet
		if !shared.IntInSlice(migrationId, migrations) {

			//Load the sql script
			content, err := os.ReadFile("migrations/" + fileName)
			if err != nil {
				log.Fatal(err)
			}

			//Execute the sql script
			log.Println("Executing migration: " + fileName)
			requests := strings.Split(string(content), ";")
			for _, request := range requests {
				if len(request) == 0 {
					continue
				}

				_, err := db.Exec(request)
				if err != nil {
					log.Fatal(err)
				}
			}
		}

	}
	log.Println("Finished migrations")
}

// getMigrationIds returns the Ids of migrations which have been applied to the database
func getMigrationIds(db *sql.DB) []int {
	var migrations []int

	//Get the current database state
	rows, err := db.Query("SELECT migrations FROM db_state")
	if err != nil {
		//Error 1146 says 'Table 'api.db_state' doesn't exist'
		//We can ignore that because we will create the table in the next step.
		if strings.Contains(err.Error(), "Error 1146") {
			return migrations
		}
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var migration int
		err = rows.Scan(&migration)
		if err != nil {
			log.Fatal(err)
		}
		migrations = append(migrations, migration)
	}

	return migrations
}

func CreateDatabaseIfNotExists(database string) {
	// connect to db using standard Go database/sql API
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("MYSQL_ROOT_USER"),
		os.Getenv("MYSQL_ROOT_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"))

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("Create database if not exists %s", database))
	if err != nil {
		if connectionString == ":@tcp(:)/" {
			log.Println("Environment variables have not been set. See https://github.com/AustrianDataLAB/IndieGameStream/tree/develop/api")
		}
		log.Fatal("Error creating database: ", err)
	}

	return
}
