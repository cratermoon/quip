package quipdb

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

//QuipRepo generates short, witty, quips from a repository
type QuipRepo struct {
	sdb *simpledb.SimpleDB
}

func listDomains(sdb *simpledb.SimpleDB) {

	params := &simpledb.ListDomainsInput{
		MaxNumberOfDomains: aws.Int64(4),
		NextToken:          aws.String("String"),
	}
	resp, err := sdb.ListDomains(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
	}

	fmt.Println(resp)
}

func countQuips(sdb *simpledb.SimpleDB) (int64, error) {
	params := &simpledb.DomainMetadataInput{
		DomainName: aws.String("quips"),
	}
	resp, err := sdb.DomainMetadata(params)
	if err != nil {
		return 0, err
	}
	return *resp.ItemCount, nil
}

func getQuip(sdb *simpledb.SimpleDB) (string, error) {

	params := &simpledb.SelectInput{
		SelectExpression: aws.String("select text from `quips`"), // Required
	}
	resp, err := sdb.Select(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return "", err
	}
	if len(resp.Items) == 0 {
		return "woops", nil
	}
	i := rand.Intn(len(resp.Items) - 1)
	return strings.TrimSpace(*resp.Items[i].Attributes[0].Value), nil
}

// Quip returns a single short, witty, quip
func (q QuipRepo) Quip() (string, error) {
	return getQuip(q.sdb)
}

func (q QuipRepo) Count() (int64, error) {
	return countQuips(q.sdb)
}

// NewQuipRepo returns a new quip repository
func NewQuipRepo() QuipRepo {
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	sdb := simpledb.New(
		session.Must(session.NewSessionWithOptions(
			session.Options{SharedConfigState: session.SharedConfigEnable}),
		),
	)

	return QuipRepo{sdb}
}
