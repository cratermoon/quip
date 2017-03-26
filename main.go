package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/svc"
)

func main() {
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	quipRepo := quipdb.NewQuipRepo()

	quip, err := quipRepo.Quip()

	if err != nil {
		fmt.Print(err)
		return
	}
	fmt.Println(quip)
	svc.Setup()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
