package quipdb

import (
	"errors"
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
		return "Experience tranquility", err
	}
	if len(resp.Items) == 0 {
		return "", errors.New("Experience tranquility")
	}
	i := rand.Intn(len(resp.Items) - 1)
	return strings.TrimSpace(*resp.Items[i].Attributes[0].Value), nil
}

// Quip returns a single short, witty, quip
func (q QuipRepo) Quip() (string, error) {
	return getQuip(q.sdb)
}

// Count returns the number of quips available in the repo
func (q QuipRepo) Count() (int64, error) {
	return countQuips(q.sdb)
}

// NewQuipRepo returns a new quip repository
func NewQuipRepo() (QuipRepo, error) {
	rand.Seed(time.Now().UnixNano() * int64(os.Getpid()))

	s, err := session.NewSessionWithOptions(session.Options{Profile: "cmdev", SharedConfigState: session.SharedConfigEnable})

	var qr QuipRepo

	if err != nil {
		return qr, err
	}

	sdb := simpledb.New(s)
	qr = QuipRepo{sdb}
	return qr, nil
}
