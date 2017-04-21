package job

import (
	"expvar"
	"log"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/twit"
)
var (
	scheduler        = expvar.NewMap("scheduler")
)

type Status struct {
	running bool
}

func (s Status) String() string {
	if s.running {
		return `"running"`
	}
	return `"stopped"`
}

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
	// idea: TakeNew() creates (and returns) a channel,
	// waits for  a little while for message on that channel
	// upon message reception, delete the quip
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
	scheduler.Add("posts", 1)
	// assuming we got here without error, tell the quipdb
	// to move the quip from newquips to the archive
}

func Schedule() {
	gocron.Every(1).Day().At("15:00").Do(post)
	s := gocron.NewScheduler()
	st := Status{true}
	scheduler.Set("status", st)
	scheduler.Add("posts", 0)
	<-s.Start()
}
