package query

import (
	"context"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
	"strings"
)

type ListHostRules struct {
	UserID   string
	HostID   string
	PageSize int
	PageNum  int
}

type ListHostRulesHandler decorator.CommandHandler[*ListHostRules, []*entity.HostRule]

func NewListHostRulesHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) ListHostRulesHandler {
	return &listHostRulesHandler{
		logger: logger,
		repos:  repos,
	}
}

type listHostRulesHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *listHostRulesHandler) Handle(ctx context.Context, query *ListHostRules) ([]*entity.HostRule, error) {
	// TODO 验证用户和主机权限

	rule, err := h.repos.HostRuleRepo.ListHostRule(ctx, &entity.ListHostRuleOptions{
		HostID:   query.HostID,
		PageSize: query.PageNum,
		PageNum:  query.PageNum,
	})
	if err != nil {
		h.logger.WithError(err).Error("failed to list host rules")
		return nil, err
	}

	var rules []*entity.HostRule
	for _, r := range rule {
		rule, err := h.repos.RuleRepo.Get(ctx, query.UserID, r.RuleID)
		if err != nil {
			h.logger.WithError(err).Error("failed to get rule")
			return nil, err
		}
		rules = append(rules, &entity.HostRule{
			Type:        rule.Type,
			CreatedAt:   ctime.FormatTimestamp(rule.CreatedAt),
			ID:          rule.ID,
			HostID:      query.HostID,
			UserID:      rule.UserID,
			Name:        rule.Name,
			Description: rule.Description,
			Port:        rule.Port,
			Proto:       rule.Proto,
			Action:      rule.Action,
			Host:        strings.Split(rule.Host, ","),
		})
	}

	h.logger.WithField("rules", rules).Info("主机规则列表")

	return rules, nil
}
