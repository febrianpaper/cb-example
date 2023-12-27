package arangodb

import "fbriansyah/client/internal/domain/senderlog"

type SenderlogModel struct {
	ID           string   `json:"_id"`
	Key          string   `json:"_key"`
	Phone        string   `json:"phone"`
	Message      string   `json:"message"`
	TemplateName string   `json:"templateName"`
	TemplateData []string `json:"templateData"`
	Status       string   `json:"status"`
}

func (m *SenderlogModel) ToDomain() senderlog.Senderlog {
	return senderlog.Senderlog{
		ID:           m.Key,
		Phone:        m.Phone,
		Message:      m.Message,
		TemplateName: m.TemplateName,
		TemplateData: m.TemplateData,
		Status:       m.Status,
	}
}

func (m *SenderlogModel) FromDomain(s senderlog.Senderlog) {
	m.Key = s.ID
	m.Phone = s.Phone
	m.Message = s.Message
	m.TemplateName = s.TemplateName
	m.TemplateData = s.TemplateData
	m.Status = s.Status
}

type CreateSenderLogModel struct {
	Phone        string   `json:"phone"`
	Message      string   `json:"message"`
	TemplateName string   `json:"templateName"`
	TemplateData []string `json:"templateData"`
	Status       string   `json:"status"`
}

func (l *CreateSenderLogModel) FromDomain(s senderlog.Senderlog) {
	l.Phone = s.Phone
	l.Message = s.Message
	l.TemplateName = s.TemplateName
	l.TemplateData = s.TemplateData
	l.Status = s.Status
}
