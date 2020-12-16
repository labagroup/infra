package twilio_test

import (
	"context"
	"testing"

	"github.com/gopub/types"
	"github.com/labagroup/infra/twilio"
)

func TestSMS_Send(t *testing.T) {
	s := twilio.NewSMS()
	err := s.Send(context.Background(), &types.PhoneNumber{
		Code:   86,
		Number: 18600366077,
	}, "【一码科技】周五会议 已取消")
	if err != nil {
		t.Errorf("Send: %v", err)
	}
}
