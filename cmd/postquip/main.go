package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"os"

	"github.com/cratermoon/quip/signing"
)

type NewQuip struct {
	Quip      string `json:"quip"`
	Signature string `json:"sig"`
}

func usage(cmd string) {
	fmt.Fprintf(os.Stderr, "usage: %s %q\n", cmd, "<quip>")
	os.Exit(1)
}

var (
	quip    = flag.String("q", "", "Provide a witty saying")
	keyFile = flag.String("k", "quip.key", "key file for posting (optional)")
	url     = flag.String("u", "http://localhost:8080/quip", "url of quip server")
	verbose = flag.Bool("v", false, "be verbose")
)

func main() {
	flag.Parse()

	if *quip == "" {
		usage(os.Args[0])
	}

	if !*verbose {
		fmt.Printf("Posting a new quip (%s) at %v\n", *quip, time.Now())
	}

	q := NewQuip{Quip: *quip}

	crt, err := ioutil.ReadFile(*keyFile)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	s := signing.Signer{Key: crt}
	q.Signature, err = s.Sign(q.Quip)
	if err != nil {
		fmt.Println(err)
		return
	}
	var j bytes.Buffer
	json.NewEncoder(&j).Encode(q)

	resp, err := http.Post(*url, "application/json", &j)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if *verbose {
		fmt.Println(string(body))
	}
}
