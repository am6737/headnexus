package query

import (
	"context"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/domain/host/entity"
)

func (h *HostHandler) Find(ctx context.Context, query *host.FindHost) ([]*host.Host, error) {
	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Filters:      query.Filters,
		Sort:         query.Sort,
		Limit:        query.Limit,
		Offset:       query.Offset,
		NetworkID:    query.NetworkID,
		IPAddress:    query.IPAddress,
		Role:         query.Role,
		Name:         query.Name,
		IsLighthouse: query.IsLighthouse,
	})
	if err != nil {
		return nil, err
	}

	var hosts []*host.Host
	for _, h := range find {
		hosts = append(hosts, host.ConvertEntityToHost(h))
	}

	return hosts, nil
}
