package job

import (
	"fmt"
	"time"
	"os"

	"github.com/jasonlvhit/gocron"
)

func post() {
	fmt.Fprintln(os.Stderr, "Yo I heard you like tasks", time.Now().Format(time.Stamp))
}

func schedule() {
	fmt.Fprintln(os.Stderr, time.Now().Format(time.Stamp))
	s := gocron.NewScheduler()
	s.Every(1).Day().At("15:00").Do(post)
	<- s.Start()
}
