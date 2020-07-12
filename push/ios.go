package push

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/gopub/environ"
	"github.com/gopub/errors"
	"github.com/gopub/log"
	"github.com/gopub/types"
	"github.com/sideshow/apns2"
	"github.com/sideshow/apns2/certificate"
)

type Notification struct {
	Token   string
	Title   string
	Message string
	Badge   int
	Env     types.Env
}

type IOSPusher struct {
	prodClient *apns2.Client
	devClient  *apns2.Client
	bundleID   string
}

func NewIOSPusher() *IOSPusher {
	p := new(IOSPusher)
	p.devClient = newClient("push.ios.dev_cert", "push.ios.dev_password").Development()
	p.prodClient = newClient("push.ios.prod_cert", "push.ios.prod_password").Production()
	p.bundleID = environ.MustString("app.bundle_id")
	return p
}

func newClient(certEnv, passEnv string) *apns2.Client {
	certFile := environ.MustString(certEnv)
	certPass := environ.String(passEnv, "")
	cert, err := certificate.FromP12File(certFile, certPass)
	if err != nil {
		log.Fatalf("Read certificate %s: %v", certFile, err)
	}
	return apns2.NewClient(cert)
}

func (p *IOSPusher) Push(ctx context.Context, n *Notification) error {
	payload := types.M{
		"aps": types.M{
			"alert": types.M{
				"title": n.Title,
				"body":  n.Message,
			},
			"badge": n.Badge,
		},
	}
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	c := p.devClient
	if n.Env == types.Prod {
		c = p.prodClient
	}
	v := &apns2.Notification{
		DeviceToken: n.Token,
		Payload:     payloadData,
		Topic:       p.bundleID,
	}
	resp, err := c.PushWithContext(ctx, v)
	if err != nil {
		return fmt.Errorf("push: %w", err)
	}
	if !resp.Sent() {
		if n.Env != types.Dev {
			p.devClient.PushWithContext(ctx, v)
		}
		return errors.Format(resp.StatusCode, resp.Reason)
	}
	return nil
}
