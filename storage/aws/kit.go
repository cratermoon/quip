// Package aws provides a toolkit of basic AWS functions
package aws

import (
	"errors"
	"io/ioutil"
	"log"
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
	defaultAWSProfile string = "quip"
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

// FileObject returns the content at keyname in the default bucket
func (k *Kit) FileObject(keyname string) ([]byte, error) {
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

// DBSelectRandom will return a randomly selected value from the results of query
func (k *Kit) DBSelectRandom(query string) (string, error) {
	params := &simpledb.SelectInput{
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

// DBCountItems returns the number of items in the given domain
func (k *Kit) DBCountItems(domain string) (int64, error) {
	params := &simpledb.DomainMetadataInput{
		DomainName: aws.String(domain),
	}
	resp, err := k.sdb.DomainMetadata(params)
	if err != nil {
		return 0, err
	}
	return *resp.ItemCount, nil
}

// DBList fetches all the attributes in a domain
func (k *Kit) DBList(attribute string, domain string) ([]string, error) {
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

func (k *Kit) DBTakeFirst(attribute string, domain string) (string, error) {
	q := fmt.Sprintf("select %s from `%s` limit 1", attribute, domain)
	params := &simpledb.SelectInput{
		SelectExpression: aws.String(q),
	}
	resp, err := k.sdb.Select(params)

	if err != nil {
		return "", err
	}
	if len(resp.Items) != 1 {
		return "", errors.New("Experience tranquility")
	}
	attrs := resp.Items[0].Attributes
	name := resp.Items[0].Name
	attr := attrs[0]
	log.Printf("Taking %s: %s", *name, *attr.Value)
	delParams := &simpledb.DeleteAttributesInput{
		DomainName: aws.String(domain),
		ItemName:   aws.String(*name),
	}
	k.sdb.DeleteAttributes(delParams)
	return *attr.Value, nil
}

// DBAdd places a new value in the domain at given attribute name
func (k *Kit) DBAdd(attribute string, value string, domain string) (string, error) {
	if domain == "" {
		domain = defaultSDBDomain
	}
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
		DomainName: aws.String(domain),
		ItemName:   aws.String(*id.Value),
	}
	// according to Amazon's documentation, the returned PutAttributesOutput
	// is an opaque struct anyway, so no point in assigning it
	_, err := k.sdb.PutAttributes(params)
	if err != nil {
		return "", err
	}
	return *id.Value, nil
}
