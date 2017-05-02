package job

import (
	"expvar"
	"log"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/twit"
)

var (
	schedvars = expvar.NewMap("scheduler")
	lastErrMsg expvar.String
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
	// TakeNew() creates (and returns) a channel,
	// waits for  a little while for message on that channel
	// upon message reception, delete the quip
	quip, c, _ := r.TakeNew()
	// if we get nothing, grab a random one from the archive
	if quip == "" {
		log.Println("Nothing new under the sun")
		quip, err = r.Quip()
	}
	t := twit.NewTwit()

	if t == nil {
		log.Println("Error creating twitter kit")
	}
	quip = quip + " #qotd"
	id, err := t.Tweet(quip)
	if err != nil {
		log.Printf("Error tweeting quip %q (%d) %s", quip, id, err)
		schedvars.Add("post-errors", 1)
		lastErrMsg.Set(err.Error())
		return
	}
	schedvars.Add("posts", 1)
	if c != nil {
		s, err := r.Add(quip)
		if err != nil {
			log.Printf("Error adding quip %s to archive %v", s, err)
			return
		}
		// assuming we got here without error, tell the quipdb
		// to cancel moving the quip back to the new list
		c <- true
	}
}

func Schedule() {
	log.Print("Scheduler starting")
	gocron.Every(1).Day().At("15:00").Do(post)
	st := Status{true}
	schedvars.Set("status", st)
	schedvars.Set("tweet-error", &lastErrMsg)
	lastErrMsg.Set("none")
	schedvars.Add("posts", 0)
	schedvars.Add("post-errors", 0)
	log.Printf("Job status: %s, posts %s", schedvars.Get("status").String(), schedvars.Get("posts").String())
	<-gocron.Start()
	st.running = false
}
