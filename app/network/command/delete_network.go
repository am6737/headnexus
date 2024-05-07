package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/network"
)

func (h *NetworkHandler) Delete(ctx context.Context, cmd *network.DeleteNetwork) error {
	h.logger.WithField("cmd", cmd).Info("delete network")

	ps, err := h.repo.UsedIPs(ctx, cmd.ID)
	if err != nil {
		h.logger.Error("failed to get used ips")
		return err
	}

	if len(ps) > 0 {
		return errors.New("network has used addresses")
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		h.logger.Error("failed to delete network")
		return err
	}

	return nil
}
