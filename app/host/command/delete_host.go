package command

import (
	"context"
	"github.com/am6737/headnexus/infra/persistence"
	pcode "github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type DeleteHost struct {
	UserID string
	ID     string
}

type DeleteHostHandler decorator.CommandHandlerNoneResponse[*DeleteHost]

func NewDeleteHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) DeleteHostHandler {
	return &deleteHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type deleteHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *deleteHostHandler) Handle(ctx context.Context, cmd *DeleteHost) error {
	h.logger.WithField("cmd", cmd).Info("Delete Request")

	host, err := h.repos.HostRepo.Get(ctx, cmd.ID)
	if err != nil {
		h.logger.WithError(err).Error("Error getting host")
		return err
	}
	if host.Owner != cmd.UserID {
		return pcode.Forbidden
	}

	if err := h.repos.NetworkRepo.ReleaseIP(ctx, host.NetworkID, host.IPAddress); err != nil {
		h.logger.WithError(err).Error("Error recycling address")
		return err
	}

	if err := h.repos.HostRepo.Delete(ctx, cmd.ID); err != nil {
		h.logger.WithError(err).Error("Error deleting host")
		return err

	}

	return nil
}
