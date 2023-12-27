package usecase

import (
	"context"
	"errors"
	"fbriansyah/client/internal/domain/senderlog"
	"fbriansyah/client/internal/domain/vendor"
	"fbriansyah/client/internal/usecase/commands"
	"fbriansyah/client/pkg/httpbreaker"
)

type (
	IUseCases interface {
		ICommands
		IQueries
	}
	ICommands interface {
		SendSms(ctx context.Context, cmd commands.SendSmsParams) error
	}
	IQueries interface{}
	UseCases struct {
		vendorRepo    vendor.Repository
		senderLogRepo senderlog.Repository
		httpClient    httpbreaker.HttpClient
		useCaseCommands
	}
	useCaseCommands struct {
		*commands.SendSmsHandler
	}
	UsecaseConfig func(*UseCases) error
)

func WithVendorRepo(repo vendor.Repository) UsecaseConfig {
	return func(u *UseCases) error {
		if repo == nil {
			return errors.New("usecase vendor repo is required")
		}
		u.vendorRepo = repo
		return nil
	}
}

func WithSenderLogRepo(repo senderlog.Repository) UsecaseConfig {
	return func(u *UseCases) error {
		if repo == nil {
			return errors.New("usecase sender log repo is required")
		}
		u.senderLogRepo = repo
		return nil
	}
}

func WithHttpClient(client httpbreaker.HttpClient) UsecaseConfig {
	return func(u *UseCases) error {
		if client == nil {
			return errors.New("usecase http client is required")
		}
		u.httpClient = client
		return nil
	}
}

func New(vendors []string, configs ...UsecaseConfig) (*UseCases, error) {
	uc := &UseCases{}

	// apply configs
	for _, config := range configs {
		err := config(uc)
		if err != nil {
			return nil, err
		}
	}
	if uc.vendorRepo == nil {
		return nil, errors.New("new usecase vendor repo is required")
	}

	if uc.senderLogRepo == nil {
		return nil, errors.New("new usecase sender log repo is required")
	}

	if uc.httpClient == nil {
		return nil, errors.New("new usecase http client is required")
	}

	// create send sms command handler
	sendSmsCommandHandler, err := commands.NewSendSmsHandler(
		vendors,
		commands.WithVendorRepo(uc.vendorRepo),
		commands.WithSenderLogRepo(uc.senderLogRepo),
		commands.WithHttpClient(uc.httpClient),
	)
	if err != nil {
		return nil, err
	}

	// register use case commands
	uc.useCaseCommands = useCaseCommands{
		SendSmsHandler: sendSmsCommandHandler,
	}

	return uc, nil
}

var _ IUseCases = (*UseCases)(nil)
