package main

import (
	"log"
	"os"
	"pricescraper-worker/internal"
	"pricescraper-worker/internal/repo"
)

func main() {
	db, err := repo.NewMongo()
	if err != nil {
		log.Fatal(err)
	}

	app, err := internal.Inject(db)
	if err != nil {
		log.Fatal(err)
	}

	port := ":" + os.Getenv("PORT")

	log.Fatal(app.Run(port))

}
