package quipdb

import (
	"github.com/cratermoon/quip/storage"
)

//QuipRepo generates short, witty, quips from a repository
type QuipRepo struct {
	kit storage.Kit
}

// Quip returns a single short, witty, quip from the archive
func (q QuipRepo) Quip() (string, error) {
	resp, err := q.kit.DBSelectRandom("select text from `quips`")

	if err != nil {
		return "Experience tranquility", err
	}

	return resp, nil
}

// TakeNew will remove and return the first quip in the new list
// and add it to the archive
func (q QuipRepo) TakeNew() (string, error) {
	resp, err := q.kit.DBTakeFirst("text", "newquips")
	if err != nil {
		return resp, err
	}
	defer q.kit.DBAdd("text", resp, "quips")
	return resp, err
}

// Count returns the number of quips available in the archive
func (q QuipRepo) Count() (int64, error) {
	return q.kit.DBCountItems("quips")
}

// List retuns the list of all quips in the archive
func (q QuipRepo) List() ([]string, error) {
	return q.kit.DBList("text", "quips")
}

// ListNew retuns the list of all new quips
func (q QuipRepo) ListNew() ([]string, error) {
	return q.kit.DBList("text", "newquips")
}

// Add will insert the given string into the newquips repo
func (q QuipRepo) Add(quip string) (string, error) {
	return q.kit.DBAdd("text", quip, "newquips")
}

// NewQuipRepo returns a new quip repository
func NewQuipRepo() (QuipRepo, error) {

	var qr QuipRepo

	kit, err := storage.NewKit()
	if err != nil {
		return qr, err
	}

	qr = QuipRepo{kit}
	return qr, nil
}
