package quipdb

import (
	"github.com/cratermoon/quip/storage"
	"log"
)

type QuipBasket interface {
	Quip() string
	Empty(bool)
}

type kitQuip struct {
	kit  storage.Kit
	quip string
}

func (q kitQuip) Quip() string {
	return q.quip
}

func (q kitQuip) Empty(success bool) {
	// just log it
	log.Printf("Empty, got success=%t", success)
}

type newQuip struct {
	kitQuip
}

func (q newQuip) Quip() string {
	return q.quip
}
func (q newQuip) Empty(success bool) {
	log.Printf("Emptying, got success=%t", success)
	if success {
		// move to archive
		q.kit.DBAdd("text", q.quip, "quips")
	} else {
		// return it to the newquips bucket
		q.kit.DBAdd("text", q.quip, "newquips")
	}
}
