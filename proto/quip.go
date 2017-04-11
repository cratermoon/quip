package proto

type QuipRequest struct{}

type QuipResponse struct {
	Quip string `json:"quip"`
	Err  string `json:"err,omitempty"`
}

type ListQuipResponse struct {
	Quips []string `json:"quips"`
	Err   string   `json:"err,omitempty"`
}

type AddQuipRequest struct {
	Quip      string `json:"quip"`
	Signature string `json:"sig"`
	UUID      string `json:"uuid,omitempty"`
}

func (a AddQuipRequest) Value() string {
	return a.UUID
}

type AddQuipResponse struct {
	Name string `json:"name"`
	Err  string `json:"err,omitempty"`
}
