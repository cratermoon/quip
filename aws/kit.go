package aws

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"strings"

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
		SelectExpression: aws.String("select text from `quips`"), // Required
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
