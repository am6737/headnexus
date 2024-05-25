package query

import (
	"context"
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/domain/host/entity"
)

func (h *HostHandler) ListHostRules(ctx context.Context, query *host.ListHostRules) ([]*entity.Rule, error) {
	// TODO 验证用户和主机权限

	rule, err := h.hostRuleRepo.ListHostRule(ctx, &entity.ListHostRuleOptions{
		HostID:   query.HostID,
		PageSize: query.PageNum,
		PageNum:  query.PageNum,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to list host rules")
		return nil, err
	}

	var rules []*entity.Rule
	for _, r := range rule {
		rule, err := h.ruleRepo.Get(ctx, query.UserID, r.RuleID)
		if err != nil {
			h.logger.WithError(err).Error("failed to get rule")
			return nil, err
		}
		rules = append(rules, rule)
	}

	h.logger.WithField("rules", rules).Info("主机规则列表")

	return rules, nil
}
