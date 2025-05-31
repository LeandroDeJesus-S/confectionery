package schemas

// Message represents a message returned by the API
type Message struct {
	Code   int      `json:"code"`
	Detail []string `json:"detail"`
}