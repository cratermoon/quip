package aws

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/simpledb"
)

// Interface to AWS
type AWSKit struct {
	session *session.Session
	sdb     *simpledb.SimpleDB
	s3      *s3.S3
}

// NewAWSKit creates a new object for AWS-related functions
func NewAWSKit() (*AWSKit, error) {

	session, err := session.NewSessionWithOptions(session.Options{Profile: "cmdev", SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, err
	}
	sdb := simpledb.New(session)
	s3 := s3.New(session)

	kit := AWSKit{
		session: session,
		sdb:     sdb,
		s3:      s3,
	}
	return &kit, nil
}

// S3Object returns the content at keyname in the default bucket
func (k *AWSKit) S3Object(keyname string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String("cmdev.com"),
		Key:    aws.String(keyname),
	}
	resp, err := k.s3.GetObject(params)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// SDBSelectRandom will return a randomly selected value from the results of query
func (k *AWSKit) SDBSelectRandom(query string) (string, error) {
	params := &simpledb.SelectInput{
		// "select text from `quips`"
		SelectExpression: aws.String(query),
	}
	resp, err := k.sdb.Select(params)

	if err != nil {
		return "Experience tranquility", err
	}
	if len(resp.Items) == 0 {
		return "", errors.New("Experience tranquility")
	}
	i := rand.Intn(len(resp.Items))
	return strings.TrimSpace(*resp.Items[i].Attributes[0].Value), nil
}

// SDBCountItems returns the number of items in the given domain
func (k *AWSKit) SDBCountItems(domain string) (int64, error) {
	params := &simpledb.DomainMetadataInput{
		DomainName: aws.String(domain),
	}
	resp, err := k.sdb.DomainMetadata(params)
	if err != nil {
		return 0, err
	}
	return *resp.ItemCount, nil
}

func (k *AWSKit) SDBList(attribute string, domain string) ([]string, error) {
	q := fmt.Sprintf("select %s from `%s`", attribute, domain)
	params := &simpledb.SelectInput{
		SelectExpression: aws.String(q),
	}
	resp, err := k.sdb.Select(params)

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

func (k *AWSKit) SDBAdd(attribute string, value string) (string, error) {
	id := &simpledb.ReplaceableAttribute{
		Name:  aws.String("id"),
		Value: aws.String(strconv.FormatInt(time.Now().Unix(), 10)),
	}
	text := &simpledb.ReplaceableAttribute{
		Name:  aws.String(attribute),
		Value: aws.String(value),
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
	_, err := k.sdb.PutAttributes(params)
	if err != nil {
		return "", err
	}
	return *id.Value, nil
}
