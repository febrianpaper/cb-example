package senderlog

type Senderlog struct {
	ID           string
	Phone        string
	Message      string
	TemplateName string
	TemplateData []string
	Status       string
}
