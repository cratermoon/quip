package models

// Profile represents a toy person
type Profile struct {
	ID        string
	Name      string
	Addresses []Location
}

// Location is a string for an address
type Location string
