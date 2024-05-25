package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/pkg/code"
)

func (h *RuleHandler) Delete(ctx context.Context, cmd *rule.DeleteRule) error {
	h.logger.WithField("cmd", cmd).Info("delete rule")

	r, err := h.repo.Get(ctx, cmd.UserID, cmd.ID)
	if err != nil {
		if errors.Is(err, code.NotFound) {
			return code.NotFound.SetMessage("rule not found")
		}
		return err
	}
	if r.UserID != cmd.UserID {
		return code.Forbidden
	}

	err = h.repo.Delete(ctx, cmd.ID)
	if err != nil {
		return err
	}

	return nil
}
