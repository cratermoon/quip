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

	"github.com/cratermoon/quip/job"
	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/server"
)

type options struct {
	keyFile string
	port    string
	verbose bool
}

var (
	Version string
	Build   string

	author        = expvar.NewString("author")
	authorContact = expvar.NewString("authorContact")
)

func readOptions() options {
	opts := options{}
	flag.StringVar(&opts.keyFile, "k", "quip.crt", "private key cert")
	flag.StringVar(&opts.port, "p", "8080", "port to serve on")
	flag.BoolVar(&opts.verbose, "v", false, "be verbose")
	flag.Parse()
	return opts
}

func main() {
	opts := readOptions()

	if opts.verbose {
		fmt.Printf("Version: %s Build: %s\n", Version, Build)
	}

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

	go job.Schedule()

	h := server.BuildServices(opts.keyFile)

	fmt.Println("Listening on port", opts.port)
	log.Fatal(http.ListenAndServe(":"+opts.port, h))
}
