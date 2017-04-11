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

	"github.com/cratermoon/quip/aws"
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

	//b, err := ioutil.ReadAll(uresp.Body)
	var uuid proto.UUIDResponse
	err = json.NewDecoder(uresp.Body).Decode(&uuid)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}
	q := proto.AddQuipRequest{Quip: *quip, UUID: uuid.UUID}

	kit, err := aws.NewKit()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	crt, err := kit.S3Object(*keyFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	s := signing.Signer{Key: crt}
	v := strings.Join([]string{q.Quip, q.UUID}, ":")
	q.Signature, err = s.Sign(v)
	if err != nil {
		fmt.Println(err)
		return
	}
	var j bytes.Buffer
	json.NewEncoder(&j).Encode(q)

	resp, err := http.Post(*url+"/quip", "application/json", &j)
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
