// Package storage defines a toolkit of basic datastore functions
package storage

import "github.com/cratermoon/quip/storage/aws"
import gcloud "github.com/cratermoon/quip/storage/gcloud"

// Kit is an interface to basic storage functions
type Kit interface {
	FileObject(name string) ([]byte, error)
	DBSelectRandom(query string) (string, error)
	DBCountItems(domain string) (int64, error)
	DBList(attribute string, domain string) ([]string, error)
	DBTakeFirst(attribute string, domain string) (string, error)
	DBAdd(attribute string, value string, domain string) (string, error)
}

// NewKit creates a new object for storage-related functions
func NewKit() (Kit, error) {
	return newAWSKit()
}
func newAWSKit() (Kit, error) {
	return aws.NewKit()
}

func mewGCloudKit() (Kit, error) {
	return gcloud.NewKit()
}
