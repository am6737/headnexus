package query

import (
	"context"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/entity"
)

type FindNetwork struct {
	Name          string
	Cidr          string
	SortDirection string
	IncludeCounts bool
	PageSize      int
	PageNumber    int
}

var DefaultQueryNetwork = &network.FindNetwork{
	SortDirection: "",
	IncludeCounts: false,
	PageSize:      10, // 默认每页大小为 10
	PageNumber:    1,  // 默认页码为 1
}

func (h *NetworkHandler) Find(ctx context.Context, query *network.FindNetwork) ([]*entity.Network, error) {
	h.logger.WithField("query", query).Info("find network")

	if query == nil {
		query = DefaultQueryNetwork
	}

	if query.SortDirection == "" {
		query.SortDirection = "asc"
	}
	if !query.IncludeCounts {
		query.IncludeCounts = true
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.PageNumber <= 0 {
		query.PageNumber = 1
	}

	networks, err := h.repo.Find(ctx, &entity.QueryNetwork{
		Name:          query.Name,
		Cidr:          query.Cidr,
		SortDirection: query.SortDirection,
		IncludeCounts: query.IncludeCounts,
		PageSize:      query.PageSize,
		PageNumber:    query.PageNumber,
	})
	if err != nil {
		h.logger.Errorf("failed to find networks: %v", err)
		return nil, err
	}

	return networks, nil
}
