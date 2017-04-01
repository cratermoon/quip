package main

import (
	"bytes"
	"encoding/json"
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

func main() {
	fmt.Println("Posting a new quip at", time.Now())
	url := "http://localhost:8080/quip"

	if len(os.Args) < 2 {
		usage(os.Args[0])
	}
	q := NewQuip{Quip: os.Args[1]}

	var err error
	q.Signature, err = signing.Sign(q.Quip)
	if err != nil {
		fmt.Println(err)
		return
	}
	var j bytes.Buffer
	json.NewEncoder(&j).Encode(q)

	resp, err := http.Post(url, "application/json", &j)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
