package cmd

import (
	"expvar"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/server"
)

var author = expvar.NewString("author")

func Serve() {
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	quipRepo, err := quipdb.NewQuipRepo()
	if err != nil {
		log.Print("error getting the quip respository: ", err)
	}

	quip, err := quipRepo.Quip()

	if err != nil {
		log.Print("error getting a quip from the repo: ", err)
	} else {
		fmt.Println(quip)
	}

	h := server.BuildServices()

	author.Set("Steven E. Newton")

	log.Fatal(http.ListenAndServe(":8080", h))
}
