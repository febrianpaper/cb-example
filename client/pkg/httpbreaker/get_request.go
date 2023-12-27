package httpbreaker

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/sony/gobreaker"
)

type GetConfig struct {
	URL     string
	Headers http.Header
}

func (h *HttpBreaker) Get(vendor string, ctx context.Context, cfg GetConfig) ([]byte, error) {
	cb, ok := vendorsCB[vendor]
	if !ok {
		return nil, ErrUnkownVendor
	}

	state := cb.State()
	if state == gobreaker.StateOpen {
		return nil, ErrVendorNotAvailable
	}

	if cfg.URL == "" {
		return nil, ErrEmptyUrl
	}

	body, err := cb.Execute(func() (interface{}, error) {
		req, err := http.NewRequest(http.MethodGet, cfg.URL, nil)
		if err != nil {
			return nil, ErrCreateNewRequest
		}

		req.Header = cfg.Headers
		req = req.WithContext(ctx)

		resp, err := h.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%s vendor=%s, method=GET, url=%s", ErrSendingRequest.Error(), vendor, cfg.URL)
		}
		defer resp.Body.Close()

		response, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, ErrReadBody
		}
		return response, nil
	})

	if err != nil {
		return nil, err
	}

	return body.([]byte), nil
}
