package twilio

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gopub/environ"
	"github.com/gopub/log"
	"github.com/gopub/types"
	"github.com/gopub/wine"
)

type SMS struct {
	phoneNumbers []string
	numIndex     int
	send         *wine.ClientEndpoint
}

func NewSMS() *SMS {
	s := &SMS{
		phoneNumbers: environ.StringSlice("twilio.numbers", nil),
		numIndex:     0,
	}
	account := environ.MustString("twilio.account")
	authToken := environ.MustString("twilio.token")
	msgURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", account)
	msgURL = environ.String("twilio.msg_url", msgURL)
	var err error
	s.send, err = wine.DefaultClient.Endpoint(http.MethodPost, msgURL)
	if err != nil {
		log.Panicf("Cannot create endpoint: %v", err)
	}
	s.send.SetBasicAuthorization(account, authToken)
	return s
}

func (s *SMS) Send(ctx context.Context, recipient *types.PhoneNumber, content string) error {
	logger := log.FromContext(ctx).With("recipient", recipient)
	s.numIndex = (s.numIndex + 1) % len(s.phoneNumbers)
	form := url.Values{}
	form.Add("To", recipient.String())
	form.Add("From", s.phoneNumbers[s.numIndex])
	form.Add("Body", content)
	var result types.M
	err := s.send.Call(ctx, form, &result)
	logger.Debug(result)
	return err
}
