package commands

import (
	"context"
	"errors"
	"fbriansyah/client/internal/domain"
	"fbriansyah/client/internal/domain/senderlog"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/internal/httpclient/halosis"
	"fbriansyah/client/internal/httpclient/infobib"
	"fbriansyah/client/pkg/httpbreaker"
	"fmt"
	"time"
)

type (
	SendSmsParams struct {
		To           string
		Message      string
		TemplateName string
		TemplateData []string
	}
	SendSmsHandler struct {
		vendorRepo       vendor.Repository
		senderLogRepo    senderlog.Repository
		httpClient       httpbreaker.HttpClient
		vendorsSenderMap map[string]domain.Sender
		vendors          []string
	}
	SendSmsConfig func(*SendSmsHandler) error
)

func WithVendorRepo(repo vendor.Repository) SendSmsConfig {
	return func(h *SendSmsHandler) error {
		if repo == nil {
			return errors.New("send_sms vendor repo is required")
		}
		h.vendorRepo = repo
		return nil
	}
}

func WithSenderLogRepo(repo senderlog.Repository) SendSmsConfig {
	return func(h *SendSmsHandler) error {
		if repo == nil {
			return errors.New("send_sms sender log repo is required")
		}
		h.senderLogRepo = repo
		return nil
	}
}

func WithHttpClient(client httpbreaker.HttpClient) SendSmsConfig {
	return func(h *SendSmsHandler) error {
		if client == nil {
			return errors.New("send_sms http client is required")
		}
		h.httpClient = client
		return nil
	}
}

// NewSendSmsHandler creates a new SendSmsHandler with the given configs.
func NewSendSmsHandler(listVendor []string, configs ...SendSmsConfig) (*SendSmsHandler, error) {
	handler := &SendSmsHandler{}
	for _, config := range configs {
		err := config(handler)
		if err != nil {
			return nil, err
		}
	}
	handler.vendors = listVendor

	// check if required fields are set
	if handler.vendorRepo == nil {
		return nil, errors.New("new send sms handler vendor repo is required")
	}

	if handler.senderLogRepo == nil {
		return nil, errors.New("new send sms handler sender log repo is required")
	}

	if handler.httpClient == nil {
		return nil, errors.New("new send sms handler http client is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	vendors, err := handler.vendorRepo.List(ctx, listVendor)
	if err != nil {
		return nil, err
	}

	// register vendors to map
	vendorsMap := make(map[string]vendor.Vendor)
	for _, vendor := range vendors {
		vendorsMap[vendor.ID] = vendor
	}

	// register vendor sender
	vendorsSenderMap := make(map[string]domain.Sender)
	if vendor, ok := vendorsMap["halosis"]; ok {
		vendorsSenderMap["halosis"] = halosis.NewSmsSender(vendor, handler.httpClient)
	}
	if vendor, ok := vendorsMap["infobib"]; ok {
		vendorsSenderMap["infobib"] = infobib.NewSmsSender(vendor, handler.httpClient)
	}
	handler.vendorsSenderMap = vendorsSenderMap

	return handler, nil
}

// SendSms sends sms to the given sender.
func (h *SendSmsHandler) SendSms(ctx context.Context, params SendSmsParams) error {

	err := h.sendSmsWithVendors(ctx, params)
	if err != nil {
		return err
	}

	return nil
}

func (h *SendSmsHandler) sendSmsWithVendors(ctx context.Context, params SendSmsParams) error {

	successSend := 0
	vendors, err := h.vendorRepo.List(ctx, h.vendors)
	if err != nil {
		return err
	}

	for _, v := range vendors {
		if !v.Setting.AllowSms {
			continue
		}

		sender := h.vendorsSenderMap[v.ID]

		err = sendSms(ctx, sender, params, &successSend, h.senderLogRepo, v)
		if err != nil {
			handleSendError(err)
			continue
		}

		if successSend > 0 {
			break
		}
	}

	if successSend == 0 {
		return errors.New("all vendor not available")
	}

	return nil
}

func sendSms(ctx context.Context, sender domain.Sender, params SendSmsParams, successSend *int, senderLogRepo senderlog.Repository, v vendor.Vendor) error {
	var log *senderlog.Senderlog
	var err error

	if v.Setting.IsSupportTemplate && params.TemplateName != "" {
		log, err = sendSmsWithTemplate(ctx, sender, params)
	} else {
		log, err = sendSmsWithoutTemplate(ctx, sender, params)
	}

	if err != nil {
		return err
	}

	errLog := senderLogRepo.Create(ctx, log)
	if errLog != nil {
		return errLog
	}

	*successSend++

	return nil
}

func handleSendError(err error) {
	if errors.Is(err, httpbreaker.ErrVendorNotAvailable) {
		// do nothing
	} else {
		// TODO: Send sender error to new relic, and send notification to slack
		fmt.Println(err)
	}
}

func sendSmsWithTemplate(ctx context.Context, sender domain.Sender, params SendSmsParams) (*senderlog.Senderlog, error) {
	// sending sms with template
	log, err := sender.SendWithTemplate(ctx, params.To, domain.TemplateConfig{
		Name: params.TemplateName,
		Data: params.TemplateData,
	})

	return log, err
}

func sendSmsWithoutTemplate(ctx context.Context, sender domain.Sender, params SendSmsParams) (*senderlog.Senderlog, error) {
	// sending sms without template
	log, err := sender.Send(ctx, params.To, params.Message)

	return log, err
}
