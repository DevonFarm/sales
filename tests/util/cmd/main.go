package main

import (
	"context"
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/DevonFarm/sales/tests/util"
)

const dbEnvVar = "TEST_DATABASE_URL"

func main() {
	godotenv.Load()
	dbSetup := func() *util.TestDB {
		connString := os.Getenv(dbEnvVar)
		if connString == "" {
			log.Fatal("TEST_DATABASE_URL is not set")
		}
		db, err := util.NewTestDB(connString)
		if err != nil {
			log.Fatal(err)
		}
		return db
	}

	// TODO: add a flag to pass in a database name
	cleanDB := flag.Bool("clean-db", false, "clean the database")
	flag.Parse()

	if *cleanDB {
		if err := dbSetup().CleanupTestDB(context.Background()); err != nil {
			log.Fatal("failed to clean database:", err)
		}
		log.Println("database cleaned")
		return
	}

	log.Println("No operation specified")
}
