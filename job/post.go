package job

import (
	"crypto/rand"
	"expvar"
	"log"
	"math"
	"math/big"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
	"github.com/cratermoon/quip/twit"
)

var (
	schedvars  = expvar.NewMap("scheduler")
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

func try(t *twit.Twit, quip string, retries int, done chan bool) {
	id, err := t.Tweet(quip + " #qotd")
	if err != nil {
		log.Printf("Error tweeting quip %q (%d) %s", quip, id, err)
		schedvars.Add("post-errors", 1)
		lastErrMsg.Set(err.Error())
		// truncated binary exponential backoff https://en.wikipedia.org/wiki/Exponential_backoff#Binary_exponential_backoff
		if retries > 10 {
			retries = 0
			return
		}
		slot := math.Pow(2, float64(retries+1))
		nBig, err := rand.Int(rand.Reader, big.NewInt(int64(slot)))
		if err != nil {
			log.Printf(`¯\_(ツ)_/¯`)
			// ¯\_(ツ)_/¯
			return
		}
		n := nBig.Int64()
		time.Sleep((time.Duration(float64(n))*time.Second)/10 + time.Microsecond)
		retries++
		try(t, quip, retries, done)
		return
	}
	if done != nil {
		done <- true
	}
}

func post() {
	log.Println("Post job running")
	var quip string
	r, err := quipdb.NewQuipRepo()
	if err != nil {
		log.Println("Post job error", err.Error())
		return
	}
	t := twit.NewTwit()

	if t == nil {
		log.Println("Error creating twitter kit")
	}

	// check newquips, ignoring errors
	// TakeNew() creates (and returns) a channel,
	// waits for  a little while for message on that channel
	// upon message reception, delete the quip
	quip, c, _ := r.TakeNew()
	defer close(c)
	// if we get nothing, grab a random one from the archive
	if quip == "" {
		log.Println("Nothing new under the sun")
		quip, err = r.Quip()
	}
	try(t, quip, 0, c)
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
