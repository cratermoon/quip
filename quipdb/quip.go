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

// Quip returns a single short, witty, quip from the archive
func (q QuipRepo) Quip() (string, error) {
	resp, err := q.kit.SDBSelectRandom("select text from `quips`")

	if err != nil {
		return "Experience tranquility", err
	}

	return resp, nil
}

// TakeNew will remove and return the first quip in the new list
// and add it to the archive
func (q QuipRepo) TakeNew() (string, error) {
	resp, err := q.kit.SDBTakeFirst("text", "newquips")
	if err != nil {
		return resp, err
	}
	defer q.kit.SDBAdd("text", resp, "quips")
	return resp, err
}

// Count returns the number of quips available in the archive
func (q QuipRepo) Count() (int64, error) {
	return q.kit.SDBCountItems("quips")
}

// List retuns the list of all quips in the archive
func (q QuipRepo) List() ([]string, error) {
	return q.kit.SDBList("text", "quips")
}

// ListNew retuns the list of all new quips
func (q QuipRepo) ListNew() ([]string, error) {
	return q.kit.SDBList("text", "newquips")
}

// Add will insert the given string into the newquips repo
func (q QuipRepo) Add(quip string) (string, error) {
	return q.kit.SDBAdd("text", quip, "newquips")
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
