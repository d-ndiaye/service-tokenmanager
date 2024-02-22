package notification

type Payload struct {
	Notifications []Notification `json:"notifications"`
}

type Notification struct {
	Topic    string   `json:"topic"`
	Message  string   `json:"message"`
	Tokens   []string `json:"tokens"`
	Platform int      `json:"platform"`
}
