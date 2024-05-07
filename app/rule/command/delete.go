package command

import (
	"context"
	"github.com/am6737/headnexus/app/rule"
)

func (h *RuleHandler) Delete(ctx context.Context, cmd *rule.DeleteRule) error {
	h.logger.WithField("cmd", cmd).Info("delete rule")
	err := h.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return nil
}
