// Package gcloud implements basic storage functions using Google Cloud services
package gcloud

// Kit is an interface to basic storage functions
type Kit struct{}

func NewKit() (*Kit, error) {
	return &Kit{}, nil
}

func (k *Kit) FileObject(name string) ([]byte, error) {
	return nil, nil
}

func (k *Kit) DBSelectRandom(query string) (string, error) {
	return "", nil
}

func (k *Kit) DBCountItems(domain string) (int64, error) {
	return 0, nil
}

func (k *Kit) DBList(attribute string, domain string) ([]string, error) {
	return nil, nil
}

func (k *Kit) DBTakeFirst(attribute string, domain string) (string, error) {
	return "", nil
}

func (k *Kit) DBAdd(attribute string, value string, domain string) (string, error) {
	return "", nil
}
