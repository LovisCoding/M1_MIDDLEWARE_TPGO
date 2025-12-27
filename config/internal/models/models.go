package models

type Mail struct {
	ID   int64  `json:"id"`
	Mail string `json:"mail"`
}

type Ressource struct {
	ID   int64   `json:"id"`
	Name string  `json:"name"`
	Type *string `json:"type,omitempty"`
}

type Alert struct {
	MailID      int64 `json:"mail_id"`
	RessourceID int64 `json:"ressource_id"`
}
