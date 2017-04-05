package signing

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Key returns the named key from Amazon S3
func Key(sess *session.Session, keyname string) ([]byte, error) {

	svc := s3.New(sess)

	params := &s3.GetObjectInput{
		Bucket: aws.String("cmdev.com"),
		Key:    aws.String(keyname),
	}
	resp, err := svc.GetObject(params)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	fmt.Println(resp)

	keyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		log.Println(err)
		return nil, err
	}

	return keyBytes, nil
}
