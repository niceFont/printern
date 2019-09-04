package internal

import (
	"encoding/json"
	"log"
	"os"
)

//WriteJSON writes the Data collected by the Crawlers to a JSON file
func WriteJSON(data []map[string]string) {

	file, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("test.json")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	f.Write(file)

	f.Sync()

}
