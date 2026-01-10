package models

import "time"

type Event struct {
	UID         string    `json:"uid"`
	Summary     string    `json:"summary"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	Start       time.Time `json:"start"`
	End         time.Time `json:"end"`
}

type EmailRequest struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Content   string `json:"content"`
}

type Resource struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Alert struct {
	MailID     int64 `json:"mail_id"`
	ResourceID int64 `json:"resource_id"`
}

type Mail struct {
	ID   int64  `json:"id"`
	Mail string `json:"mail"`
}
