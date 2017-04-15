package job

import (
	"log"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
)

func post() {
	log.Println("Post job running")
	var quip string
	r, err := quipdb.NewQuipRepo()
	if err != nil {
		log.Println("Post job error", err.Error())
		return
	}
	// check newquips, ignoring errors
	quip, _ = r.TakeNew()
	// if we get nothing, grab a random one from the archive
	if quip == "" {
		log.Println("Nothing new under the sun")
		quip, err = r.Quip()
	} else {
		r.Add(quip)
	}

	log.Println("Posting quip", quip)
}

func schedule() {
	gocron.Every(1).Day().At("15:00").Do(post)
	s := gocron.NewScheduler()
	<-s.Start()
}
