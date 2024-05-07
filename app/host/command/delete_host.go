package command

import (
	"context"
	"github.com/am6737/headnexus/app/host"
)

func (h *HostHandler) Delete(ctx context.Context, cmd *host.DeleteHost) error {
	h.logger.WithField("cmd", cmd).Info("Delete Request")

	host, err := h.repo.Get(ctx, cmd.ID)
	if err != nil {
		h.logger.WithError(err).Error("Error getting host")
		return err
	}

	if err := h.nch.RecycleAddress(ctx, host.NetworkID, host.IPAddress); err != nil {
		h.logger.WithError(err).Error("Error recycling address")
		return err
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		h.logger.WithError(err).Error("Error deleting host")
		return err
	}

	return nil
}
