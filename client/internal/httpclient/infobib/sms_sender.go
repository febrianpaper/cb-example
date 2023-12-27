package infobib

import (
	"context"
	"encoding/json"
	"errors"
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/senderlog"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/pkg/httpbreaker"
	"net/http"
)

type SmsSender struct {
	vendor vendor.Vendor
	client httpbreaker.HttpClient
}

func NewSmsSender(vendor vendor.Vendor, client httpbreaker.HttpClient) *SmsSender {
	return &SmsSender{
		vendor: vendor,
		client: client,
	}
}

func (ss *SmsSender) Send(ctx context.Context, to string, message string) (*senderlog.Senderlog, error) {
	url := ss.vendor.Setting.GetSmsUrl()
	req := newSmsRequest(to, message)
	requestBody, err := json.Marshal(req)
	if err != nil {
		return nil, errors.New("failed to marshal request body: vendor=infobib, method:send")
	}

	body, err := ss.client.Post(ss.vendor.Name, ctx, httpbreaker.PostConfig{
		URL: url,
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: requestBody,
	})
	if err != nil {
		return nil, err
	}

	var response SmsResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("error unmarshaling sms response: vendor=infobib, method:send")
	}

	log := &senderlog.Senderlog{
		ID:           response.Messages[0].MessageId,
		Phone:        response.Messages[0].To,
		Message:      message,
		TemplateName: "",
		TemplateData: []string{},
		Status:       response.Messages[0].Status.Name,
	}

	return log, nil
}

func (*SmsSender) SendWithTemplate(ctx context.Context, to string, cfg domain.TemplateConfig) (*senderlog.Senderlog, error) {
	return nil, errors.New("not implemented: vendor=infobib, method:sendWithTemplate")
}

var _ domain.Sender = (*SmsSender)(nil)
