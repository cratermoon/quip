package storage

import (
	"errors"

	"github.com/cratermoon/quip/models"
	"github.com/cratermoon/quip/uuid"
)

// ProfileStorage holds Profiles
type ProfileStorage struct {
	store map[string]models.Profile
}

// Add stores a Profile. Will overwrite any entry
// with the same ID
func (s *ProfileStorage) Add(name string, address string) (models.Profile, error) {
	l := models.Location(address)
	p := models.Profile{Name: name, Addresses: []models.Location{l}}
	id, err := uuid.NewUUID()
	if err != nil {
		return p, err
	}
	p.ID = id
	s.store[p.ID] = p
	return p, nil
}

// Get returns the Profile with the given ID
func (s *ProfileStorage) Get(id string) (models.Profile, error) {
	p, ok := s.store[id]
	if !ok {
		return p, errors.New("profile not found")
	}
	return p, nil
}

// NewProfileStorage creates and returns an empty ProfileStorage
func NewProfileStorage() ProfileStorage {
	return ProfileStorage{make(map[string]models.Profile)}
}
