package models

type Message struct {
	ID            string            `json:"id" bson:"id"`
	Sender        string            `json:"sender" bson:"sender"`
	UserID        string            `json:"user_id" bson:"user_id"`
	BccRecipients []string          `json:"bcc_recipients"`
	CcRecipients  []string          `json:"cc_recipients"`
	Recipients    []string          `json:"recipients"`
	Target        string            `json:"target" bson:"target"`
	Title         string            `json:"title" bson:"title"`
	Body          string            `json:"body" bson:"body"`
	TemplateID    string            `json:"template_id" bson:"template_id"`
	DataMap       map[string]string `json:"data_map" bson:"data_map"`
	Ts            int64             `json:"ts" bson:"ts"`
}
