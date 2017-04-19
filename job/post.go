package job

import (
	"log"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/twit"
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
	// we should probably just check to see
	// if there are new ones, and defer actually removing
	// a quip from the new list until we've successfully
	// posted it
	quip, _ = r.TakeNew()
	// if we get nothing, grab a random one from the archive
	if quip == "" {
		log.Println("Nothing new under the sun")
		quip, err = r.Quip()
	} else {
		// don't add it to the archive until we are done
		defer r.Add(quip)
	}
	t := twit.NewTwit()

	if t == nil {
		log.Println("Error creating twitter kit")
	}
	t.Tweet(quip)
}

func schedule() {
	gocron.Every(1).Day().At("15:00").Do(post)
	s := gocron.NewScheduler()
	<-s.Start()
}
