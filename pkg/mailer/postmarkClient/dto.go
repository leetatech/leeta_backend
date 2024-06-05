package postmarkClient

type EmailWithTemplateResponse struct {
	To          string `json:"To"`
	SubmittedAt string `json:"SubmittedAt"`
	MessageID   string `json:"MessageID"`
	ErrorCode   int    `json:"ErrorCode"`
	Message     string `json:"Message"`
}

type ErrorResponse struct {
	ErrorCode int    `json:"ErrorCode"`
	Message   string `json:"Message"`
}
