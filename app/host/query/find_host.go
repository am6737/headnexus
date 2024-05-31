package query

import (
	"context"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type FindHost struct {
	Limit        int
	Offset       int
	IsLighthouse bool
	UserID       string
	NetworkID    string
	IPAddress    string
	Role         string
	Name         string
	Filters      map[string]interface{}
	Sort         map[string]interface{}
}

type FindHostHandler decorator.CommandHandler[*FindHost, []*Host]

func NewFindHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) FindHostHandler {
	return &findHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type findHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *findHostHandler) Handle(ctx context.Context, query *FindHost) ([]*Host, error) {

	if query.Role == "" {
		query.Role = "none"
	}

	find, err := h.repos.HostRepo.Find(ctx, &entity.HostFindOptions{
		UserID:    query.UserID,
		Filters:   query.Filters,
		Sort:      query.Sort,
		Limit:     query.Limit,
		Offset:    query.Offset,
		NetworkID: query.NetworkID,
		IPAddress: query.IPAddress,
		Role:      query.Role,
		Name:      query.Name,
	})
	if err != nil {
		return nil, err
	}

	var hosts []*Host
	for _, h := range find {
		hosts = append(hosts, ConvertEntityToHost(h))
	}

	return hosts, nil
}

func ConvertEntityToHost(h *entity.Host) *Host {
	var online bool
	if h.Status == entity.Online {
		online = true
	}
	return &Host{
		ID:         h.ID,
		Name:       h.Name,
		IPAddress:  h.IPAddress,
		PublicIP:   h.PublicIP,
		Port:       h.Port,
		Online:     online,
		Role:       h.Role,
		LastSeenAt: ctime.FormatTimeSince(h.LastSeenAt),
		NetworkID:  h.NetworkID,
		Tags:       h.Tags,
		CreatedAt:  ctime.FormatTimestamp(h.CreatedAt),
	}
}
