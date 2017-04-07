// Package aws provides a toolkit of basic AWS functions
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

// Kit is an interface to basic AWS functions
type Kit struct {
	session *session.Session
	sdb     *simpledb.SimpleDB
	s3      *s3.S3
}

const (
	defaultAWSProfile string = "cmdev"
	defaultS3Bucket   string = "cmdev.com"
	defaultSDBDomain  string = "newquips"
)

// NewKit creates a new object for AWS-related functions
func NewKit() (*Kit, error) {

	session, err := session.NewSessionWithOptions(session.Options{Profile: defaultAWSProfile, SharedConfigState: session.SharedConfigEnable})
	if err != nil {
		return nil, err
	}
	sdb := simpledb.New(session)
	s3 := s3.New(session)

	kit := Kit{
		session: session,
		sdb:     sdb,
		s3:      s3,
	}
	return &kit, nil
}

// S3Object returns the content at keyname in the default bucket
func (k *Kit) S3Object(keyname string) ([]byte, error) {
	params := &s3.GetObjectInput{
		Bucket: aws.String(defaultS3Bucket),
		Key:    aws.String(keyname),
	}
	resp, err := k.s3.GetObject(params)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(resp.Body)
}

// SDBSelectRandom will return a randomly selected value from the results of query
func (k *Kit) SDBSelectRandom(query string) (string, error) {
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
func (k *Kit) SDBCountItems(domain string) (int64, error) {
	params := &simpledb.DomainMetadataInput{
		DomainName: aws.String(domain),
	}
	resp, err := k.sdb.DomainMetadata(params)
	if err != nil {
		return 0, err
	}
	return *resp.ItemCount, nil
}

// SDBList fetches all the attributes in a domain
func (k *Kit) SDBList(attribute string, domain string) ([]string, error) {
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

// SDBAdd places a new value in the domain at given attribute name
func (k *Kit) SDBAdd(attribute string, value string) (string, error) {
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
		DomainName: aws.String(defaultSDBDomain),
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
