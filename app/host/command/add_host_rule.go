package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/domain/host/entity"
)

func (h *HostHandler) AddHostRule(ctx context.Context, cmd *host.AddHostRule) ([]*entity.Rule, error) {
	h.logger.WithField("cmd", cmd).Info("add host rule")
	if len(cmd.Rules) == 0 {
		return nil, errors.New("no rule ids provided")
	}

	// TODO 判断用户权限

	r, err := h.repo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.Error("error getting host")
		return nil, err
	}

	if err := h.hostRuleRepo.AddHostRule(ctx, r.ID, cmd.Rules...); err != nil {
		h.logger.WithError(err).Error("add host rule failed")
		return nil, err
	}

	var rules []*entity.Rule
	for _, ruleID := range cmd.Rules {
		rule, err := h.repos.RuleRepo.Get(ctx, cmd.UserID, ruleID)
		if err != nil {
			h.logger.WithError(err).Error("get host rule failed")
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
