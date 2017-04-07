package quipdb

import (
	"math/rand"
	"time"

	"os"

	"github.com/cratermoon/quip/aws"
)

//QuipRepo generates short, witty, quips from a repository
type QuipRepo struct {
	kit *aws.Kit
}

// Quip returns a single short, witty, quip
func (q QuipRepo) Quip() (string, error) {
	resp, err := q.kit.SDBSelectRandom("select text from `quips`")

	if err != nil {
		return "Experience tranquility", err
	}

	return resp, nil
}

// Count returns the number of quips available in the repo
func (q QuipRepo) Count() (int64, error) {
	return q.kit.SDBCountItems("quips")
}

// List retuns the list of all quips
func (q QuipRepo) List() ([]string, error) {
	return q.kit.SDBList("text", "quips")
}

// Add will insert the given string into the quip repo
func (q QuipRepo) Add(quip string) (string, error) {
	return q.kit.SDBAdd("text", quip)
}

// NewQuipRepo returns a new quip repository
func NewQuipRepo() (QuipRepo, error) {
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	var qr QuipRepo

	kit, err := aws.NewKit()
	if err != nil {
		return qr, err
	}

	qr = QuipRepo{kit}
	return qr, nil
}
