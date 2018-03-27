package plug

// Message is the standard response sent to user with the generic routes set by the model methods
type Message struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}
