package senderlog

import "context"

type Repository interface {
	Create(ctx context.Context, log *Senderlog) error
}
