package command

import (
	"context"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/repository"
	"github.com/sirupsen/logrus"
)

var _ network.CommandHandler = &NetworkHandler{}

func NewNetworkHandler(repo repository.NetworkRepository, logger *logrus.Logger) *NetworkHandler {
	return &NetworkHandler{
		logger: logger,
		repo:   repo,
	}
}

type NetworkHandler struct {
	logger *logrus.Logger

	repo repository.NetworkRepository
}

func (h *NetworkHandler) AllocateAutoAddress(ctx context.Context, networkID string) (string, error) {
	ip, err := h.repo.AutoAllocateIP(ctx, networkID)
	if err != nil {
		return "", err
	}
	return ip, nil
}

func (h *NetworkHandler) AllocateStaticAddress(ctx context.Context, networkID string, addr string) (string, error) {
	err := h.repo.AllocateIP(ctx, networkID, addr)
	if err != nil {
		return "", err
	}
	return addr, err
}

func (h *NetworkHandler) RecycleAddress(ctx context.Context, networkID string, addr string) error {
	if err := h.repo.ReleaseIP(ctx, networkID, addr); err != nil {
		return err
	}

	return nil
}

func (h *NetworkHandler) UsedAddresses(ctx context.Context, networkID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}

func (h *NetworkHandler) AvailableAddresses(ctx context.Context, networkID string) ([]string, error) {
	//TODO implement me
	panic("implement me")
}
