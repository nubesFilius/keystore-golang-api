package main

import (
	"io/ioutil"
	"keystore-golang-api/api"
	"log"
	"net/http"
)

func main() {
	// check precondition: keys directory is in place
	_, err := ioutil.ReadDir("./keys")
	if err != nil {
		log.Fatalf("Could not read keys directory. %s", err)
	}

	log.Println("[INFO] Listening on http://localhost:80")
	err = http.ListenAndServe(":80", api.Router())
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}
