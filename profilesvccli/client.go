package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cratermoon/quip/models"
	"github.com/cratermoon/quip/profilesvccli/naming"
)

type ProfileResponse struct {
	Status string         `json:"status"`
	Person models.Profile `json:"person"`
}

type Person struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func main() {

	rand.Seed(time.Now().Unix() * int64(os.Getpid()))

	url := "http://localhost:8080/profile"

	name := fmt.Sprintf("%s %s", naming.RandomGivenname(), naming.RandomSurname())

	address := fmt.Sprintf("%d %s %s", naming.RandomStreetNumber(), naming.RandomColor(), naming.RandomMoniker())

	person := Person{Name: name, Address: address}
	var b []byte
	buf := bytes.NewBuffer(b)
	err := json.NewEncoder(buf).Encode(person)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.Post(url, "application/json", buf)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var pr ProfileResponse
	err = json.NewDecoder(resp.Body).Decode(&pr)
	if err != nil {
		log.Fatal("Decode: ", err)
	}
	fmt.Printf("ID: %s\nName: %s\nAddresses: %s\n", pr.Person.ID, pr.Person.Name, pr.Person.Addresses)

}
