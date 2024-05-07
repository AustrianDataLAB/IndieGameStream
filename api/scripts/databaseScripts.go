package scripts

import (
	"api/shared"
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
	"strconv"
	"strings"
)

func ConnectToDatabase() *sql.DB {
	// connect to db using standard Go database/sql API
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		viper.GetString("DATABASE.USER"),
		viper.GetString("DATABASE.PASSWORD"),
		viper.GetString("DATABASE.HOST"),
		viper.GetString("DATABASE.PORT"),
		viper.GetString("DATABASE.NAME"))

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

	//For each migration script
	for _, file := range files {
		//If the file is not a .sql file, ignore it
		if file.Name()[len(file.Name())-4:] != ".sql" {
			continue
		}

		//get its id
		migrationId, err := strconv.Atoi(strings.Split(file.Name(), "_")[0])
		if err != nil {
			log.Fatal(err)
		}
		//If the migration has not been applied yet
		if !shared.IntInSlice(migrationId, migrations) {

			//Load the sql script
			content, err := os.ReadFile("migrations/" + file.Name())
			if err != nil {
				log.Fatal(err)
			}

			//Execute the sql script
			log.Println("Executing migration: " + file.Name())
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

// getMigrationIds returns the Ids of migrations which have been applied to the database.
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
		viper.GetString("DATABASE.USER"),
		viper.GetString("DATABASE.PASSWORD"),
		viper.GetString("DATABASE.HOST"),
		viper.GetString("DATABASE.PORT"))

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("Create database if not exists %s", database))
	if err != nil {
		log.Fatal(err)
	}

	return
}
