package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// KeyMetadata represents the format in which we will cache and store keys
type KeyMetadata struct {
	KeyPem       []byte    `json:"key_pem"`
	ID           string    `json:"key_id"`
	InvalidAfter time.Time `json:"expires"`
}

type GetKeyListOutput struct {
	KeyIDList []string `json:"keys,omitempty"`
}

func readKeys(w http.ResponseWriter, r *http.Request) {
	//get files in keys directory
	files, err := ioutil.ReadDir("./keys")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not get keys")
		return
	}
	//initialize return type
	jsonArray := GetKeyListOutput{}
	//append filenames to the return array
	for _, f := range files {
		jsonArray.KeyIDList = append(jsonArray.KeyIDList, f.Name())
	}
	//marshall and send it
	bytes, err := json.Marshal(jsonArray)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not marshall keys")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(bytes))
	return
}

func writeKey(w http.ResponseWriter, r *http.Request) {
	//read the request
	reqBytes, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not read request")
		return
	}
	//marshall the request body onto a KeyMetadata
	var keyMeta KeyMetadata
	err = json.Unmarshal(reqBytes, &keyMeta)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Request was not of the right type")
		return
	}
	//create a file with the keyID as the filename and write the request bytes to it
	err = ioutil.WriteFile("keys/"+keyMeta.ID, reqBytes, 0777)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "Could not create new file")
		return
	}
	//success
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, fmt.Sprintf("Stored key %s", keyMeta.ID))
	log.Printf("Stored key %s", keyMeta.ID)
	return
}

func readKey(w http.ResponseWriter, r *http.Request) {
	//get the key id from the url
	vars := mux.Vars(r)
	keyID := vars["key_id"]

	if keyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "No key id specified")
		return
	}
	//read the file for the key id specified
	fileData, err := ioutil.ReadFile("./keys/" + keyID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "File for key specified not found")
		return
	}
	//send the bytes if found
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(fileData))
	return
}
