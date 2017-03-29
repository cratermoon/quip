package models

// Profile represents a toy person
type Profile struct {
	ID        string     `json:"id"`
	Name      string     `json:"name"`
	Addresses []Location `json:"addresses"`
}

// Location is a string for an address
type Location string
