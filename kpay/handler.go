package kpay

import context "golang.org/x/net/context"

// Server describes the KPay server structure
type Server struct{}

// RegisterPayment register a new payment given the payment input
func (s *Server) RegisterPayment(ctx context.Context, in *RegisterPaymentInput) (*RegisterPaymentOutput, error) {
	return &RegisterPaymentOutput{}, nil
}
