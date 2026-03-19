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

	cleanDB := flag.String("clean-db", "", "clean the specified database")
	flag.Parse()

	if *cleanDB != "" {
		if err := dbSetup().CleanupTestDB(context.Background(), *cleanDB); err != nil {
			log.Fatalf("failed to clean database %s: %v", *cleanDB, err)
		}
		log.Println("database cleaned")
		return
	}

	log.Println("No operation specified")
}
