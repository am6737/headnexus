package query

import (
	"context"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/entity"
)

func (h *NetworkHandler) Get(ctx context.Context, query *network.GetNetwork) (*entity.Network, error) {
	get, err := h.repo.Get(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	return get, nil
}
