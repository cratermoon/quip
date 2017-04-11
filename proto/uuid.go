package proto

type UUIDResponse struct {
	Status string `json"status,omitempty"`
	UUID   string `json:"uuid"`
}

func (u UUIDResponse) Value() string {
	return u.UUID
}
