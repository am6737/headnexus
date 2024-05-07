package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/segmentio/ksuid"
)

type CreateHostResponse struct {
	ID              string
	Name            string
	NetworkID       string
	Role            string
	CreatedAt       int64
	LastSeenAt      int64
	StaticAddresses []string
	IPAddress       string
	Port            int
	IsLighthouse    bool
	Tags            map[string]interface{}
}

func (h *HostHandler) Create(ctx context.Context, cmd *host.CreateHost) (*host.Host, error) {
	h.logger.WithField("cmd", cmd).Info("Create Request")

	if cmd.IsLighthouse && len(cmd.StaticAddresses) == 0 {
		return nil, errors.New("lighthouse主机必须指定静态ip")
	}

	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Name: cmd.Name,
	})
	if err != nil {
		return nil, err
	}

	if len(find) > 0 {
		return nil, errors.New("host name already exists")
	}

	if cmd.IPAddress != "" {
		_, err := h.nch.AllocateStaticAddress(ctx, cmd.NetworkID, cmd.IPAddress)
		if err != nil {
			h.logger.WithError(err).Error("allocate static address failed")
			return nil, err
		}
	} else {
		ip, err := h.nch.AllocateAutoAddress(ctx, cmd.NetworkID)
		if err != nil {
			h.logger.WithError(err).Error("allocate auto address failed")
			return nil, err
		}
		cmd.IPAddress = ip
		h.logger.WithField("ip", ip).
			WithField("hostname", cmd.Name).
			WithField("network_id", cmd.NetworkID).
			Info("allocate auto address")
	}

	id := ksuid.New().String()
	cfg := config.GenerateConfigTemplate()
	cfg.Tun.IP = cmd.IPAddress

	e := &entity.Host{
		ID:              id,
		Name:            cmd.Name,
		NetworkID:       cmd.NetworkID,
		Role:            cmd.Role,
		StaticAddresses: cmd.StaticAddresses,
		IPAddress:       cmd.IPAddress,
		Port:            cmd.Port,
		IsLighthouse:    cmd.IsLighthouse,
		Tags:            cmd.Tags,
		Config:          cfg,
	}

	h.logger.WithField("entity", e).Info("create host")

	r, err := h.repo.Create(ctx, e)
	if err != nil {
		if err := h.nch.RecycleAddress(ctx, cmd.NetworkID, cmd.IPAddress); err != nil {
			h.logger.WithError(err).Error("recycle address failed")
		}
		return nil, err
	}

	r2 := host.ConvertEntityToHost(r)
	return r2, nil
}
