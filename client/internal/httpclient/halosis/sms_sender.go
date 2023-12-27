package halosis

import (
	"context"
	"encoding/json"
	"errors"
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/senderlog"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/pkg/httpbreaker"
	"fmt"
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

	requestBody := fmt.Sprintf(`{"to":"%s","message":"%s"}`, to, message)

	body, err := ss.client.Post(ss.vendor.Name, ctx, httpbreaker.PostConfig{
		URL: url,
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: []byte(requestBody),
	})
	if err != nil {
		return nil, err
	}

	var response SmsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("error unmarshaling sms response")
	}

	log := &senderlog.Senderlog{
		ID:           fmt.Sprintf("%d", response.MSGID),
		Phone:        to,
		Message:      message,
		TemplateName: "",
		TemplateData: []string{},
		Status:       response.Status,
	}

	return log, err
}

func (ss *SmsSender) SendWithTemplate(ctx context.Context, to string, cfg domain.TemplateConfig) (*senderlog.Senderlog, error) {
	requestBody := SmsWithTemplateRequest{
		To: to,
		Template: SmsTemplateSetting{
			Name:       cfg.Name,
			Components: make([]SmsComponent, len(cfg.Data)),
		},
	}

	for i, data := range cfg.Data {
		requestBody.Template.Components[i] = SmsComponent{
			Text: data,
		}
	}

	requestBodyJson, err := json.Marshal(requestBody)
	if err != nil {
		return nil, errors.New("error marshaling sms request")
	}

	body, err := ss.client.Post(ss.vendor.Name, ctx, httpbreaker.PostConfig{
		URL: ss.vendor.Setting.GetSmsUrl(),
		Headers: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: requestBodyJson,
	})

	if err != nil {
		return nil, err
	}

	var response SmsResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, errors.New("error unmarshaling sms response")
	}

	log := &senderlog.Senderlog{
		ID:           fmt.Sprintf("%d", response.MSGID),
		Phone:        to,
		Message:      "",
		TemplateName: cfg.Name,
		TemplateData: cfg.Data,
		Status:       response.Status,
	}

	return log, nil
}

var _ domain.Sender = (*SmsSender)(nil)
