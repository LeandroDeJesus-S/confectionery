package schemas

// Message represents a message returned by the API
type Message struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}