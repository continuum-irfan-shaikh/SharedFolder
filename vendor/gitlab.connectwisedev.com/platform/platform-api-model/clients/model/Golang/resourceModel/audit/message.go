package audit

// Message - A message template used to store
// and update messages in the Audit service
type Message struct {
	Object string `json:"object,omitempty"`
	Code   string `json:"code,omitempty"`
	Value  string `json:"value,omitempty"`
}
