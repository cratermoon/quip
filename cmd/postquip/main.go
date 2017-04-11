package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"os"

	"github.com/cratermoon/quip/proto"
	"github.com/cratermoon/quip/signing"
)

func usage(cmd string) {
	fmt.Fprintf(os.Stderr, "usage: %s %q\n", cmd, "<quip>")
	os.Exit(1)
}

var (
	quip    = flag.String("q", "", "Provide a witty saying")
	keyFile = flag.String("k", "quip.key", "key file for posting (optional)")
	url     = flag.String("u", "http://localhost:8080", "url of quip server")
	verbose = flag.Bool("v", false, "be verbose")
)

func check(e error) {
	if e != nil {
		fmt.Fprintln(os.Stderr, e.Error())
		os.Exit(2)
	}
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *quip == "" {
		flag.PrintDefaults()
		return
	}

	if !*verbose {
		fmt.Printf("Posting a new quip (%s) at %v\n", *quip, time.Now())
	}

	uresp, err := http.Get(*url + "/uuid")
	check(err)

	var uuid proto.UUIDResponse
	err = json.NewDecoder(uresp.Body).Decode(&uuid)
	check(err)

	q := proto.AddQuipRequest{Quip: *quip, UUID: uuid.UUID}

	crt, err := ioutil.ReadFile(*keyFile)
	check(err)

	s := signing.Signer{Key: crt}
	v := strings.Join([]string{q.Quip, q.UUID}, ":")
	q.Signature, err = s.Sign(v)
	check(err)

	var j bytes.Buffer
	json.NewEncoder(&j).Encode(q)

	resp, err := http.Post(*url+"/quip", "application/json", &j)
	check(err)
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	if *verbose {
		fmt.Println(string(body))
	}
}
