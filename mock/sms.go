package mock

import (
	"context"

	"github.com/gopub/log"
	"github.com/gopub/types"
)

type SMS struct {
}

func (s *SMS) Send(ctx context.Context, recipient *types.PhoneNumber, content string) error {
	log.FromContext(ctx).Debugf("%v %s", recipient, content)
	return nil
}
