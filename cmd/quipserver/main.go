package main

import (
	"expvar"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/server"
)

var (
	Version string
	Build   string

	keyFile       = flag.String("k", "quip.crt", "private key cert")
	port          = flag.String("p", "8080", "port to serve on")
	verbose       = flag.Bool("v", false, "be verbose")
	author        = expvar.NewString("author")
	authorContact = expvar.NewString("authorContact")
)

func main() {
	fmt.Printf("Version: %s Build: %s\n", Version, Build)
	flag.Parse()

	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	quipRepo, err := quipdb.NewQuipRepo()
	if err != nil {
		log.Fatal("error getting the quip respository: ", err)
	}

	quip, err := quipRepo.Quip()

	if err != nil {
		log.Fatal("error getting a quip from the repo: ", err)
	} else {
		log.Print(quip)
	}

	author.Set("Steven E. Newton")
	authorContact.Set("snewton@treetopllc.com")

	h := server.BuildServices(*keyFile)

	fmt.Println("Listening on port", *port)
	log.Fatal(http.ListenAndServe(":"+*port, h))
}
