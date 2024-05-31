package command

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	net2 "github.com/am6737/headnexus/pkg/net"
	"github.com/sirupsen/logrus"
)

type UpdateHost struct {
	UserID    string
	HostID    string
	Name      string
	IPAddress string
	PublicIP  string
	Port      int
}

type UpdateHostHandler decorator.CommandHandler[*UpdateHost, *entity.Host]

func NewUpdateHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) UpdateHostHandler {
	return &updateHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type updateHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *updateHostHandler) Handle(ctx context.Context, cmd *UpdateHost) (*entity.Host, error) {
	host, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get host")
		return nil, err
	}
	if host.Owner != cmd.UserID {
		return nil, code.Forbidden
	}

	if cmd.Name != "" {
		host.Name = cmd.Name
	}
	if cmd.IPAddress != "" {
		network, err := h.repos.NetworkRepo.Get(ctx, host.NetworkID)
		if err != nil {
			return nil, err
		}
		if !net2.IsIPInCIDR(cmd.IPAddress, network.Cidr) {
			return nil, code.InvalidParameter.SetMessage("Ip地址不在网络cidr中")
		}

		if err := h.repos.NetworkRepo.ReleaseIP(ctx, network.ID, host.IPAddress); err != nil {
			h.logger.WithError(err).Error("failed to release ip")
			return nil, err
		}

		if err := h.repos.NetworkRepo.AllocateIP(ctx, network.ID, cmd.IPAddress); err != nil {
			h.logger.WithError(err).Error("failed to allocate ip")
			return nil, err
		}

		host.IPAddress = cmd.IPAddress
	}
	if cmd.PublicIP != "" {
		host.PublicIP = cmd.PublicIP
	}
	if cmd.Port != 0 {
		host.Port = cmd.Port
	}

	update, err := h.repos.HostRepo.Update(ctx, host)
	if err != nil {
		h.logger.WithError(err).Error("failed to update host")
		return nil, err
	}

	return update, nil
}
