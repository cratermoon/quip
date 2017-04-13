package job

import (
	"fmt"

	"github.com/jasonlvhit/gocron"

	"github.com/cratermoon/quip/quipdb"
)

func post() {
	fmt.Println("I am runnning task.")
	var quip string
	r, err := quipdb.NewQuipRepo()
	if err != nil {
		return
	}
	// check newquips
	// if newquips not empty
	//   get quip
        //   move quip to archive
        // else
	quip, err = r.Quip()
	if err != nil {
		return
	}
	fmt.Println("posting quip", quip)
}


func schedule() {
	gocron.Every(1).Day().At("15:00").Do(post)
	s := gocron.NewScheduler()
	<- s.Start()
}
