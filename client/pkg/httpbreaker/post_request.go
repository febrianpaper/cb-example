package httpbreaker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/sony/gobreaker"
)

type PostConfig struct {
	URL     string
	Headers http.Header
	Body    []byte
}

func (h *HttpBreaker) Post(vendor string, ctx context.Context, cfg PostConfig) ([]byte, error) {
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
		req, err := http.NewRequest(http.MethodPost, cfg.URL, bytes.NewBuffer(cfg.Body))
		if err != nil {
			return nil, ErrCreateNewRequest
		}
		req.Header = cfg.Headers
		req = req.WithContext(ctx)

		resp, err := h.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("%s vendor=%s, method=POST, url=%s", ErrSendingRequest.Error(), vendor, cfg.URL)
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
