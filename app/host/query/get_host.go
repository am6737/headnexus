package query

import (
	"context"
	"github.com/am6737/headnexus/app/host"
)

func (h *HostHandler) Get(ctx context.Context, query *host.GetHost) (*host.Host, error) {
	r1, err := h.repo.Get(ctx, query.ID)
	if err != nil {
		h.logger.Errorf("error getting host: %v", err)
		return nil, err
	}

	r2 := host.ConvertEntityToHost(r1)

	return r2, nil
}
