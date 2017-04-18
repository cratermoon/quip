// Package gcloud implements basic storage functions using Google Cloud services
package gcloud


// StorageKit is an interface to basic storage functions
type Kit struct {
}

	FileObject(name string) ([]byte, error)
	DBSelectRandom(query string) (string, error)
	DBCountItems(domain string) (int64, error)
	DBList(attribute string, domain string) ([]string, error)
	DBTakeFirst(attribute string, domain string) (string, error)
	DBAdd(attribute string, value string, domain string) (string, error)
