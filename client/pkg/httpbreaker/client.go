package httpbreaker

import (
	"context"
	"time"

	"net/http"

	"github.com/sony/gobreaker"
)

var (
	vendorsCB map[string]*gobreaker.CircuitBreaker
)

type HttpClient interface {
	Get(vendor string, ctx context.Context, cfg GetConfig) ([]byte, error)
	Post(vendor string, ctx context.Context, cfg PostConfig) ([]byte, error)
}

type HttpBreaker struct {
	client *http.Client
}

type Settings struct {
	// Interval is the cyclic period of the closed state for the CircuitBreaker to clear the internal Counts.
	// If Interval is less than or equal to 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
	Interval time.Duration
	// Timeout is the period of the open state, after which the state of the CircuitBreaker becomes half-open.
	// If Timeout is less than or equal to 0, the timeout value of the CircuitBreaker is set to 60 seconds.
	Timeout time.Duration
	// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
	// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state. If ReadyToTrip is nil, default ReadyToTrip is used. Default ReadyToTrip returns true when the number of consecutive failures is more than 5.
	ReadyToTrip func(counts gobreaker.Counts) bool
}

func New(vendors []string, client *http.Client, settings *Settings) *HttpBreaker {
	vendorsCB = make(map[string]*gobreaker.CircuitBreaker)
	for _, vendor := range vendors {
		vendorsCB[vendor] = gobreaker.NewCircuitBreaker(gobreaker.Settings{
			Name:        vendor,
			Interval:    settings.Interval,
			Timeout:     settings.Timeout,
			ReadyToTrip: settings.ReadyToTrip,
		})
	}

	return &HttpBreaker{
		client: client,
	}
}

var _ HttpClient = (*HttpBreaker)(nil)
