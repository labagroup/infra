package twilio_test

import (
	"context"
	"testing"

	"github.com/gopub/types"
	"github.com/libnat/infra/twilio"
)

func TestSMS_Send(t *testing.T) {
	s := twilio.NewSMS()
	err := s.Send(context.Background(), &types.PhplanNumber{
		Code:   0,
		Number: 0,
	}, "Hello")
	if err != nil {
		t.Errorf("Send: %v", err)
	}
}
