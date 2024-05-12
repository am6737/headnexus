package command

import (
	"context"
	"github.com/am6737/headnexus/app/user"
)

func (h *UserHandler) Delete(ctx context.Context, cmd *user.DeleteUser) error {
	h.logger.WithField("cmd", cmd).Info("delete rule")
	err := h.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return nil
}
