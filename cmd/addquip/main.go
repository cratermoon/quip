package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/cratermoon/quip/storage/aws"
)

const (
	attributeName string = "text"
	maxQuipLength int    = 274
)

var (
	quip    = flag.String("q", "", "Provide a witty saying")
	verbose = flag.Bool("v", false, "be verbose")
	dryrun  = flag.Bool("d", false, "Dry run - show what would be done")
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if *quip == "" {
		flag.PrintDefaults()
		return
	}

	if len(*quip) > maxQuipLength {
		fmt.Fprintf(os.Stderr, "Maximum quip length (%d) exceeded, got %d\n",
			maxQuipLength, len(*quip))
		return
	}

	if *verbose {
		fmt.Printf("Posting a new quip (%s) at %v\n", *quip, time.Now().Format(time.Stamp))
	}

	kit, err := aws.NewKit()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	if *dryrun {
		fmt.Printf("Adding new quip '%s' as attribute 'text'\n", *quip)
		return
	}
	id, err := kit.DBAdd("text", *quip, "newquips")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	if !*verbose {
		fmt.Printf("Quip %s posted\n", id)
	}

}
