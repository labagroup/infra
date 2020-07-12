package twilio

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gopub/environ"
	"github.com/gopub/errors"
	"github.com/gopub/log"
	"github.com/gopub/types"
)

type SMS struct {
	Account       string
	AuthToken     string
	msgURL        string
	phplanNumbers []string
	numIndex      int
}

func NewSMS() *SMS {
	s := &SMS{
		Account:       environ.String("twilio.account", ""),
		AuthToken:     environ.String("twilio.token", ""),
		phplanNumbers: environ.StringSlice("twilio.numbers", nil),
		numIndex:      0,
	}
	if len(s.phplanNumbers) == 0 {
		log.Fatalf("missing twilio.numbers")
	}
	defaultMsgURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", s.Account)
	s.msgURL = environ.String("twilio.msg_url", defaultMsgURL)
	return s
}

func (s *SMS) Send(ctx context.Context, recipient *types.PhoneNumber, content string) error {
	logger := log.FromContext(ctx).With("recipient", recipient)
	s.numIndex = (s.numIndex + 1) % len(s.phplanNumbers)
	form := &url.Values{}
	form.Add("To", recipient.String())
	form.Add("From", s.phplanNumbers[s.numIndex])
	form.Add("Body", content)
	req, err := http.NewRequest("POST", s.msgURL, strings.NewReader(form.Encode()))
	if err != nil {
		logger.Error(err)
		return nil
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(s.Account+":"+s.AuthToken)))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Error(err)
		return nil
	}
	defer resp.Body.Close()
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error(err)
		return nil
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
		logger.Infof("OK: %s", content)
		return nil
	}

	logger.Error(string(respBody))
	var result types.M
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return errors.Format(resp.StatusCode, err.Error())
	}
	return errors.Format(resp.StatusCode, result.String("message"))
}
