package quipdb

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"os"

	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

//QuipRepo generates short, witty, quips from a repository
type QuipRepo struct {
	sdb *simpledb.SimpleDB
}

// Quip returns a single short, witty, quip
func (q QuipRepo) Quip() (string, error) {
	params := &simpledb.SelectInput{
		SelectExpression: aws.String("select text from `quips`"), // Required
	}
	resp, err := q.sdb.Select(params)

	if err != nil {
		return "Experience tranquility", err
	}
	if len(resp.Items) == 0 {
		return "", errors.New("Experience tranquility")
	}
	i := rand.Intn(len(resp.Items))
	return strings.TrimSpace(*resp.Items[i].Attributes[0].Value), nil
}

// Count returns the number of quips available in the repo
func (q QuipRepo) Count() (int64, error) {
	params := &simpledb.DomainMetadataInput{
		DomainName: aws.String("quips"),
	}
	resp, err := q.sdb.DomainMetadata(params)
	if err != nil {
		return 0, err
	}
	return *resp.ItemCount, nil
}

// List retuns the list of all quips
func (q QuipRepo) List() ([]string, error) {
	params := &simpledb.SelectInput{
		SelectExpression: aws.String("select text from `quips`"), // Required
	}
	resp, err := q.sdb.Select(params)

	if err != nil {
		return nil, err
	}
	count := len(resp.Items)
	if count == 0 {
		return nil, errors.New("Experience tranquility")
	}

	quips := make([]string, count, count)
	for i, item := range resp.Items {
		quips[i] = *item.Attributes[0].Value
	}
	return quips, nil
}

// Add will insert the given string into the quip repo
func (q QuipRepo) Add(quip string) (string, error) {
	id := &simpledb.ReplaceableAttribute{
		Name:  aws.String("id"),
		Value: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	text := &simpledb.ReplaceableAttribute{
		Name:  aws.String("text"),
		Value: aws.String(quip),
	}
	attributes := make([]*simpledb.ReplaceableAttribute, 2, 2)
	attributes[0] = id
	attributes[1] = text
	params := &simpledb.PutAttributesInput{
		Attributes: attributes,
		DomainName: aws.String("newquips"),
		ItemName:   aws.String(*id.Value),
	}
	// according to Amazon's documentation, the returned PutAttributesOutput
	// is an opaque struct anyway
	_, err := q.sdb.PutAttributes(params)
	if err != nil {
		return "", err
	}
	return *id.Value, nil
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
