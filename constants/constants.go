package constants

import "time"

const (
	MAIL_SERVICE_SMTP    string = "smtp"
	MAIL_SERVICE_MAILGUN string = "mailgun"
	MAIL_SERVICE_SES     string = "ses"

	AUTH_NONE  string = "none"
	AUTH_PLAIN string = "plain"
	AUTH_LOGIN string = "login"
)

const MAIL_SERVER_URI string = "http://localhost:8025/api/v1"

type MailpitMessage struct {
	Attachments int `json:"Attachments"`
	Bcc         []struct {
		Address string `json:"Address"`
		Name    string `json:"Name"`
	} `json:"Bcc"`
	Cc []struct {
		Address string `json:"Address"`
		Name    string `json:"Name"`
	} `json:"Cc"`
	Created time.Time `json:"Created"`
	From    struct {
		Address string `json:"Address"`
		Name    string `json:"Name"`
	} `json:"From"`
	ID        string `json:"ID"`
	MessageID string `json:"MessageID"`
	Read      bool   `json:"Read"`
	ReplyTo   []struct {
		Address string `json:"Address"`
		Name    string `json:"Name"`
	} `json:"ReplyTo"`
	Size    int      `json:"Size"`
	Snippet string   `json:"Snippet"`
	Subject string   `json:"Subject"`
	Tags    []string `json:"Tags"`
	To      []struct {
		Address string `json:"Address"`
		Name    string `json:"Name"`
	} `json:"To"`
}

type MailpitMessagesResponse struct {
	Messages []MailpitMessage `json:"messages"`
}
