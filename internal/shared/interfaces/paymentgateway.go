package interfaces

import (
	"context"
)

type IPaymentGateway interface {
	CreateOrder(ctx context.Context, amount int, currency string) (providerID string, providerName string)
}
