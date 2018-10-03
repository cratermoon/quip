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

func post() {
	log.Println("Post job running")
	r, err := quipdb.NewQuipRepo()
	if err != nil {
		log.Println("Post job error", err.Error())
		return
	}
	// check newquips, ignoring errors
	// TakeNew() creates (and returns) a QuipBasket,
	quip, _ := r.TakeNew()
	// if we get nothing, grab a random one from the archive
	if quip == nil {
		log.Println("Nothing new under the sun")
		quip, err = r.Quip()
	}

	t := twit.NewTwit()
	if t == nil {
		log.Println("Error creating twitter kit")
		quip.Empty(false)
		return
	}
	tweet := quip.Quip() + " #qotd"
	retries := 0
	success := false;
	for retries < 10 && !success {
		id, err := t.Tweet(tweet)
		if err != nil {
			log.Printf("Error tweeting quip %q (%d) %s", tweet, id, err)
			schedvars.Add("post-errors", 1)
			lastErrMsg.Set(err.Error())
			// truncated binary exponential backoff https://en.wikipedia.org/wiki/Exponential_backoff#Binary_exponential_backoff
			slot := math.Pow(2, float64(retries))
			nBig, err := rand.Int(rand.Reader, big.NewInt(int64(slot)))
			if err != nil {
				// there's really no recovering here. Give up
				// ¯\_(ツ)_/¯
				log.Printf("%s: %v",`¯\_(ツ)_/¯`, err)
			}
			n := nBig.Int64()
			retries++;
			time.Sleep((time.Duration(n) * time.Second) / 10)
		} else {
			success = true
		}
	}
	if success {
		schedvars.Add("posts", 1)
	}
	quip.Empty(success)
	return
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
