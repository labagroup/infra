package mock

import (
	"context"

	"github.com/gopub/log"
	"github.com/labagroup/infra/push"
)

type Pusher struct {
}

func (p *Pusher) Push(ctx context.Context, n *push.Notification) error {
	log.FromContext(ctx).Debugf("%+v", n)
	return nil
}
